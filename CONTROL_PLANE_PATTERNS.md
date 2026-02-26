# Cilium Control Plane Engineering Patterns

A curated guide to the most clever and production-hardened engineering patterns found
in Cilium's control plane. Each pattern includes the motivation, the design, and
references to the source code for deeper study.

---

## Table of Contents

1. [Hive: Dependency Injection & Lifecycle Framework](#1-hive-dependency-injection--lifecycle-framework)
2. [StateDB: In-Process Immutable Radix Tree Database](#2-statedb-in-process-immutable-radix-tree-database)
3. [Reconciler: Desired-State Synchronization Engine](#3-reconciler-desired-state-synchronization-engine)
4. [Resource: Kubernetes Event Streaming Abstraction](#4-resource-kubernetes-event-streaming-abstraction)
5. [SelectorCache: Namespace-Indexed Identity Matching](#5-selectorcache-namespace-indexed-identity-matching)
6. [Trigger: Event Folding with Reason Deduplication](#6-trigger-event-folding-with-reason-deduplication)
7. [Controller: ExitReason Error Differentiation](#7-controller-exitreason-error-differentiation)
8. [Lock Patterns: Sortable Mutexes & Semaphored RWLock](#8-lock-patterns-sortable-mutexes--semaphored-rwlock)
9. [Promise: Generic Futures with Context Support](#9-promise-generic-futures-with-context-support)
10. [Backoff: Cluster-Size Adaptive Retry](#10-backoff-cluster-size-adaptive-retry)
11. [etcd Lease Manager: Reference-Counted Pooling](#11-etcd-lease-manager-reference-counted-pooling)
12. [Identity Allocation: Withheld Set for Smooth Restart](#12-identity-allocation-withheld-set-for-smooth-restart)
13. [BPF Object Cache: Content-Addressed Compilation](#13-bpf-object-cache-content-addressed-compilation)
14. [Set Container: Single-Element Optimization](#14-set-container-single-element-optimization)
15. [StoppableWaitGroup: Graceful Shutdown Primitive](#15-stoppablewaitgroup-graceful-shutdown-primitive)
16. [OnDemand: Reference-Counted Lazy Resource Lifecycle](#16-ondemand-reference-counted-lazy-resource-lifecycle)

---

## 1. Hive: Dependency Injection & Lifecycle Framework

**Location:** `pkg/hive/`

**Problem:** Large distributed systems like Cilium have hundreds of components with
complex interdependencies. Managing startup order, shutdown order, configuration, and
testability becomes intractable with manual wiring.

**Solution:** Hive is a modular framework built on `go.uber.org/dig` that provides:
- **Cells** as composable building blocks (modules, providers, invokers)
- **Automatic lifecycle management** with LIFO shutdown ordering
- **Module decorators** that auto-scope loggers, health, and job groups per module

**Key Pattern — Module Decorators (auto-scoping):**

Every module automatically receives its own scoped logger, database handle, and job
group without any explicit wiring:

```go
// pkg/hive/hive.go
moduleDecorators := []cell.ModuleDecorator{
    // Each module gets a logger tagged with its module ID
    func(mid cell.ModuleID) logging.FieldLogger {
        return logging.DefaultSlogLogger.With(logfields.LogSubsys, string(mid))
    },
    // Each module gets a scoped database handle
    func(db *statedb.DB, mid cell.ModuleID) *statedb.DB {
        return db.NewHandle(string(mid))
    },
}
```

**Key Pattern — LIFO Lifecycle with Rollback:**

Startup hooks execute sequentially. If any hook fails, only the hooks that
successfully started are stopped — in reverse order:

```go
// Only numStarted hooks get stopped, not the full list
func (lc *DefaultLifecycle) Stop(log *slog.Logger, ctx context.Context) error {
    var errs error
    for ; lc.numStarted > 0; lc.numStarted-- {
        hook := lc.hooks[lc.numStarted-1]  // Reverse order
        if err := hook.Stop(ctx); err != nil {
            errs = errors.Join(errs, err)
        }
    }
    return errs
}
```

**Key Pattern — In/Out Structs for Clean Dependencies:**

```go
type providerParams struct {
    cell.In
    Lifecycle cell.Lifecycle
    Config    *Config
    Logger    *slog.Logger
}

type output struct {
    cell.Out
    Provider   types.Provider
    HealthInfo string
}

func NewProvider(params providerParams) (output, error) {
    // Tagged struct fields replace long parameter lists
}
```

**Why it's clever:** Zero boilerplate wiring. A single `hive.New()` call auto-includes
jobs, health, metrics, and StateDB. Module ID validation uses compile-time regex
enforcement. The framework catches misuse at startup with panics rather than allowing
subtle runtime failures.

---

## 2. StateDB: In-Process Immutable Radix Tree Database

**Location:** `vendor/github.com/cilium/statedb/` (vendored from `github.com/cilium/statedb`)

**Problem:** Control plane state needs to be shared across many goroutines with
transactional consistency, efficient change tracking, and zero-copy reads — all
without external database overhead.

**Solution:** StateDB is a fully in-process database built on persistent (immutable)
adaptive radix trees. It provides MVCC-style concurrency where readers never block
writers.

**Key Pattern — Lock-Free Read Transactions:**

Read transactions are just atomic pointer loads — zero contention:

```go
type readTxn []*tableEntry

func (r *readTxn) getTableEntry(meta TableMeta) *tableEntry {
    return (*r)[meta.tablePos()]  // O(1) immutable snapshot access
}
```

**Key Pattern — Two-Phase Atomic Commit:**

Write transactions commit in three phases to minimize lock duration:

```go
// Phase 1: Build new immutable tree nodes (no lock contention)
for pos := range txn.tableEntries {
    table.indexes[i], txn = idx.commit()  // Creates immutable nodes
}

// Phase 2: Atomic root swap (brief lock)
db.mu.Lock()
db.root.Store(&root)  // All readers now see new snapshot
db.mu.Unlock()

// Phase 3: Notify watchers (no root lock held)
for _, txn := range txnToNotify {
    txn.notify()  // Close watch channels
}
```

**Key Pattern — Graveyard for Delete Tracking:**

Deleted objects are moved to a "graveyard" so that observers can see deletions.
A background GC worker uses low-watermark tracking across all delete trackers
to clean up:

```go
// Find the most-behind observer
lowWatermark := table.revision
for _, dt := range deleteTrackers {
    rev := dt.getRevision()
    if rev < lowWatermark {
        lowWatermark = rev
    }
}
// Delete graveyard entries below the low watermark
```

**Key Pattern — Sorted Mutex Acquisition (Deadlock-Free):**

Multi-table write transactions acquire table locks in deterministic sorted order:

```go
wtxn := db.WriteTxn(table1, table2, table3)  // Locks sorted by sequence number
```

**Why it's clever:** Structural sharing means write transactions only copy the
changed path in the radix tree — unchanged subtrees are shared between snapshots.
Watch channels provide reactive notifications without polling. The LPM (Longest
Prefix Match) index variant enables efficient IP routing table lookups.

---

## 3. Reconciler: Desired-State Synchronization Engine

**Location:** `vendor/github.com/cilium/statedb/reconciler/` (vendored from `github.com/cilium/statedb`)

**Problem:** The control plane stores desired state in StateDB, but that state must
be synchronized to external systems (BPF maps, kernel routes, etc.) with proper
error handling, retries, and periodic full reconciliation.

**Solution:** A generic reconciler framework that watches StateDB tables and
drives idempotent operations to synchronize desired state with actual state.

**Key Pattern — Per-Object Status with Atomic ID:**

Each object carries a status with an atomic ID to prevent stale updates from
multiple reconcilers:

```go
type Status struct {
    Kind      StatusKind  // Pending, Done, Error, Refreshing
    UpdatedAt time.Time
    Error     *string
    ID        uint64      // Prevents duplicate reconciliation
}
```

**Key Pattern — Exponential Backoff Priority Queue:**

Failed objects are placed in a priority queue ordered by retry time, with
exponential backoff:

```go
type retries struct {
    queue    *retryPrioQueue  // Ordered by retry time
    revQueue *retryPrioQueue  // Ordered by revision (progress tracking)
    backoff  exponentialBackoff // 100ms → 200ms → 400ms → ... → 1min
}
```

**Key Pattern — Three Reconciliation Modes:**

1. **Incremental:** Processes only changed objects (via StateDB change iterator)
2. **Full (Prune):** Periodically diffs desired state against actual state
3. **Refresh:** Marks old objects as pending at a controlled rate

```go
// The reconcile loop selects based on triggers
select {
case <-tableWatchChan:          // Table changed → incremental
case <-pruneTickerChan:         // Timer → full reconcile
case <-r.retries.Wait():        // Retry timer → retry failed
case <-r.externalPruneTrigger:  // Manual → force full
}
```

**Why it's clever:** The framework inverts control — you implement only the
target-system operations (`Update`, `Delete`, `Prune`), and the framework handles
the state machine, batching, retries, and metrics. Optional `BatchOperations`
allow efficient bulk synchronization.

---

## 4. Resource: Kubernetes Event Streaming Abstraction

**Location:** `pkg/k8s/resource/`

**Problem:** Kubernetes informers are low-level, error-prone (forgetting to handle
tombstones, race conditions on cache sync), and hard to test.

**Solution:** A higher-level abstraction that provides two access modes —
event streaming and synchronized store — with comprehensive edge case handling.

**Key Pattern — Lazy Initialization with Promise:**

Resources don't start informers until first accessed:

```go
r.storeResolver, r.storePromise = promise.New[Store[T]]()

func (r *resource[T]) startWhenNeeded() {
    <-r.needed  // Block until Events() or Store() called
    store, informer := r.newInformer()
    r.storeResolver.Resolve(&typedStore[T]{store})  // Fulfill promise
    informer.Run(r.ctx.Done())
}
```

**Key Pattern — Panic-on-Forgotten-Done() Safety:**

Uses Go runtime finalizers to detect forgotten `Done()` calls:

```go
doneFinalizer := func(done *bool) {
    panic(fmt.Sprintf(
        "%s has a broken event handler that did not call Done()",
        s.debugInfo))
}

event.Done = func(err error) {
    runtime.SetFinalizer(eventDoneSentinel, nil)  // Clear on success
    // ... handle error ...
}
runtime.SetFinalizer(eventDoneSentinel, doneFinalizer)
```

**Key Pattern — Synthetic Delete Events:**

Tracks last-known object state to emit proper delete events even when the
object has already vanished from the store:

```go
obj, exists, err := store.GetByKey(workItem.key)
if !exists {
    deletedObject, ok := lastKnownObjects.Load(workItem.key)
    if ok {
        event.Kind = Delete
        event.Object = deletedObject  // Emit with last known state
    }
}
```

**Key Pattern — UID-Aware Deletion:**

Handles rapid object recreations by tracking UIDs:

```go
func (l *lastKnownObjects[T]) DeleteByUID(key Key, objToDelete T) {
    if obj, ok := l.objs[key]; ok {
        if getUID(obj) == getUID(objToDelete) {  // UID must match
            delete(l.objs, key)
        }
    }
}
```

**Why it's clever:** The three-phase event stream (initial replay → sync marker →
incremental updates) provides a clean abstraction over Kubernetes informer semantics.
The finalizer-based safety net catches bugs that would otherwise cause silent data
loss in production.

---

## 5. SelectorCache: Namespace-Indexed Identity Matching

**Location:** `pkg/policy/selectorcache.go`

**Problem:** Network policy selectors must be matched against potentially thousands
of security identities. Naive O(n) scanning on every identity change is too slow.

**Solution:** A namespace-indexed cache that narrows the search space, combined with
async notifications to avoid lock contention.

**Key Pattern — Namespace Index for Fast Matching:**

```go
type scIdentityCache struct {
    ids         map[identity.NumericIdentity]*scIdentity
    byNamespace map[string]map[*scIdentity]struct{}  // Fast lookup
}

func (c *scIdentityCache) selections(sel *identitySelector) iter.Seq[identity.NumericIdentity] {
    return func(yield func(id identity.NumericIdentity) bool) {
        namespaces := sel.source.SelectedNamespaces()
        if len(namespaces) > 0 {
            // Only iterate identities in selected namespaces
            for _, ns := range namespaces {
                for id := range c.byNamespace[ns] {
                    if sel.source.Matches(id.lbls) {
                        if !yield(id.NID) { return }
                    }
                }
            }
        } else {
            // Full scan only when no namespace filter
            for nid, id := range c.ids { /* ... */ }
        }
    }
}
```

**Key Pattern — Async FIFO Notifications:**

Identity updates notify policy users asynchronously to avoid holding the cache
lock during expensive endpoint regenerations:

```go
func (sc *SelectorCache) handleUserNotifications() {
    for {
        sc.userMutex.Lock()
        for len(sc.userNotes) == 0 {
            sc.userCond.Wait()
        }
        notifications := sc.userNotes  // Snapshot
        sc.userNotes = nil
        sc.userMutex.Unlock()

        // Process outside lock — callbacks can't block cache
        for _, n := range notifications {
            n.user.IdentitySelectionUpdated(n.selector, n.added, n.deleted)
            n.wg.Done()
        }
    }
}
```

**Key Pattern — Transactional Selector Updates:**

Multiple identity changes batch into a single write transaction, then commit
atomically with a revision bump:

```go
func (sc *SelectorCache) commit() {
    sc.revision++
    sc.readableSelections = sc.writeableSelections.Commit()
    readTxn := types.GetSelectorSnapshot(sc.readableSelections, sc.revision)
    sc.readTxn.Store(&readTxn)  // Atomic publish to all readers
}
```

**Why it's clever:** The namespace index transforms O(n) identity matching into
O(k) where k is the number of identities in the relevant namespace. The lazy
notification handler starts only on first use (`sync.Once`).

---

## 6. Trigger: Event Folding with Reason Deduplication

**Location:** `pkg/trigger/trigger.go`

**Problem:** Rapid successive events (e.g., many identity allocations) would
cause expensive operations (e.g., checkpointing) to execute redundantly.

**Solution:** A trigger mechanism that folds multiple events into a single
execution, with minimum interval enforcement and reason tracking.

**Key Pattern — Non-Blocking Signal with Reason Tracking:**

```go
func (t *Trigger) TriggerWithReason(reason string) {
    t.mutex.Lock()
    t.trigger = true
    if t.numFolds == 0 {
        t.waitStart = time.Now()  // Latency tracking from first event
    }
    t.numFolds++
    t.foldedReasons.add(reason)  // Deduplicate reasons
    t.mutex.Unlock()

    select {
    case t.wakeupChan <- struct{}{}:
    default:  // Non-blocking — channel already signaled
    }
}
```

**Key Pattern — MinInterval Backpressure:**

The waiter enforces a minimum interval between executions, preventing thrashing:

```go
func (t *Trigger) waiter() {
    if delayNeeded, delay := t.needsDelay(); delayNeeded {
        time.Sleep(delay)
    }
    // Execute with folded metrics
    t.params.MetricsObserver.PostRun(callDuration, callLatency, numFolds)
}
```

**Why it's clever:** Reason deduplication helps debug which events triggered work.
The fold count metric reveals how much work was saved. The latency metric (time
from first trigger to execution) reveals backpressure impact.

---

## 7. Controller: ExitReason Error Differentiation

**Location:** `pkg/controller/controller.go`

**Problem:** Background controllers need to distinguish between "exit cleanly"
and "error, please retry" without resorting to string parsing or sentinel values.

**Solution:** A custom `ExitReason` type that embeds the error interface but is
handled differently by the controller loop:

```go
type ExitReason struct {
    error  // Embeds error but isn't treated as one
}

// In the run loop:
err = params.DoFunc(params.Context)
if err != nil {
    var exitReason ExitReason
    if errors.As(err, &exitReason) {
        // Clean exit — no error count increment
        c.recordSuccess(params.Health)
        interval = time.Duration(math.MaxInt64)  // Wait forever
    } else {
        // Real error — linear backoff
        errorRetries++
        interval = time.Duration(errorRetries) * params.ErrorRetryBaseDuration
    }
}
```

**Why it's clever:** Uses Go's `errors.As` type assertion for clean control flow.
The controller still records the reason for observability but treats it as success.

---

## 8. Lock Patterns: Sortable Mutexes & Semaphored RWLock

**Location:** `pkg/lock/`

**Key Pattern — Conditional Deadlock Detection via Build Tags:**

Zero overhead in production; full detection in debug builds:

```go
// lock_debug.go (lockdebug build tag):
func (i *internalMutex) Unlock() {
    if sec := time.Since(i.Time).Seconds(); sec >= 0.1 {
        printStackTo(sec, debug.Stack(), os.Stderr)  // Alert: held >100ms
    }
    i.Mutex.Unlock()
}
```

**Key Pattern — Sortable Mutexes (Deadlock-Free Multi-Lock):**

Automatically sorts mutexes by sequence number before locking:

```go
func (s SortableMutexes) Lock() {
    sort.Sort(s)  // Deterministic order → no circular wait
    for _, mu := range s {
        mu.Lock()
    }
}
```

**Key Pattern — Write-to-Read Lock Downgrade:**

Uses a weighted semaphore to implement RWLock with graceful downgrade:

```go
const maxReaders = 1 << 30  // 1 billion concurrent readers

func (i *SemaphoredMutex) Lock() {
    i.semaphore.Acquire(context.Background(), maxReaders)  // Exclusive
}

func (i *SemaphoredMutex) RLock() {
    i.semaphore.Acquire(context.Background(), 1)  // Shared
}

// Downgrade: exclusive → shared without releasing
func (i *SemaphoredMutex) UnlockToRLock() {
    i.semaphore.Release(maxReaders - 1)  // Release N-1, keep 1
}
```

**Why it's clever:** The semaphored mutex supports context cancellation (unlike
`sync.RWMutex`) and enables atomic lock downgrade — critical for patterns where
you compute under a write lock then want to hold a read lock while notifying.

---

## 9. Promise: Generic Futures with Context Support

**Location:** `pkg/promise/`

**Problem:** Go lacks built-in futures/promises. Components need to wait for
asynchronous initialization results with proper cancellation support.

**Solution:** A generic promise implementation with separate `Resolver` and
`Promise` interfaces for controlled visibility:

```go
func (p *promise[T]) Await(ctx context.Context) (value T, err error) {
    cleanupCancellation := context.AfterFunc(ctx, func() {
        p.Lock()
        defer p.Unlock()
        p.cond.Broadcast()  // Wake up Wait() on context cancel
    })
    defer cleanupCancellation()

    p.Lock()
    defer p.Unlock()
    for p.state == promiseUnresolved && ctx.Err() == nil {
        p.cond.Wait()
    }
    // Return value, rejection error, or context error
}
```

**Why it's clever:** Uses `context.AfterFunc` to bridge context cancellation with
`sync.Cond` — a notoriously tricky integration. Supports functional composition
via `Map()` for transforming promised values.

---

## 10. Backoff: Cluster-Size Adaptive Retry

**Location:** `pkg/backoff/backoff.go`

**Problem:** Fixed backoff intervals cause thundering herd in large clusters.
After a control plane restart, hundreds of nodes retry simultaneously.

**Solution:** Logarithmic scaling based on cluster size:

```go
func ClusterSizeDependantInterval(baseInterval time.Duration, numNodes int) time.Duration {
    waitNanoseconds := float64(baseInterval.Nanoseconds()) * math.Log1p(float64(numNodes))
    return time.Duration(int64(waitNanoseconds))
}
// 1 node:    ~41s      256 nodes:  ~5m32s
// 4 nodes:   ~1m36s    1024 nodes: ~6m55s
```

**Key Pattern — Auto-Reset on Idle:**

```go
func (b *Exponential) Wait(ctx context.Context) error {
    if b.ResetAfter > b.Max {
        if time.Since(b.lastBackoffStart) > b.ResetAfter {
            b.Reset()  // Recover from stale backoff state
        }
    }
    b.attempt++
    // ...
}
```

**Why it's clever:** `math.Log1p` provides sub-linear growth — doubling cluster
size only adds a constant to the backoff. Auto-reset prevents permanent stalling
after transient failures resolve.

---

## 11. etcd Lease Manager: Reference-Counted Pooling

**Location:** `pkg/kvstore/etcd_lease.go`

**Problem:** etcd has limited lease capacity. Creating one lease per key exhausts
the limit in large clusters.

**Solution:** Reference-counted lease pooling with thundering herd prevention:

```go
func (elm *etcdLeaseManager) GetSession(ctx context.Context, key string) (*concurrency.Session, error) {
    // 1. Already assigned? Reuse.
    if leaseID := elm.keys[key]; leaseID != client.NoLease {
        return elm.leases[leaseID].session, nil
    }

    // 2. Current lease not full? Add to it.
    if info := elm.leases[elm.current]; info != nil && info.count < elm.limit {
        info.count++
        elm.keys[key] = elm.current
        return info.session, nil
    }

    // 3. Find any non-full lease.
    for lease, info := range elm.leases {
        if info.count < elm.limit { /* reuse */ }
    }

    // 4. All full — acquire new lease with thundering herd prevention.
    acquiring := elm.acquiring
    if acquiring == nil {
        elm.acquiring = make(chan struct{})  // Signal channel
    } else {
        <-acquiring  // Wait for concurrent acquisition
        return elm.GetSession(ctx, key)  // Retry
    }
}
```

**Why it's clever:** The four-tier lookup (assigned → current → any → new)
minimizes lease creation. The `acquiring` channel prevents multiple goroutines
from creating leases simultaneously.

---

## 12. Identity Allocation: Withheld Set for Smooth Restart

**Location:** `pkg/identity/cache/local.go`

**Problem:** When the Cilium agent restarts, previously allocated identity IDs
should be reused to avoid unnecessary policy churn and endpoint regeneration.

**Solution:** A "withheld" set of recently-used IDs that are reserved for reuse:

```go
func (l *localIdentityCache) getNextFreeNumericIdentity(
    idCandidate identity.NumericIdentity,
) (identity.NumericIdentity, error) {
    // Try to reuse the old identity first
    if _, taken := l.identitiesByID[idCandidate]; !taken {
        return idCandidate, nil  // Smooth migration
    }

    // Search for next free, skipping withheld IDs
    for {
        _, taken := l.identitiesByID[candidate]
        _, withheld := l.withheldIdentities[candidate]
        if !taken && !withheld {
            return candidate, nil
        }
        // On wrap-around, release withheld IDs if space is tight
    }
}
```

**Why it's clever:** Withheld IDs act as a reservation system — they're preferentially
given back to the same labels on restart but released under pressure. This balances
identity stability with capacity.

---

## 13. BPF Object Cache: Content-Addressed Compilation

**Location:** `pkg/datapath/loader/`

**Problem:** BPF programs must be compiled for each endpoint, but compilation is
expensive. Many endpoints share the same datapath configuration.

**Solution:** Content-addressed caching using configuration hashing:

```go
type objectCache struct {
    baseHash datapathHash
    objects  map[string]*cachedSpec
}

func (o *objectCache) UpdateDatapathHash(nodeCfg *LocalNodeConfiguration) error {
    newHash := hashDatapath(o.ConfigWriter, nodeCfg)
    if bytes.Equal(newHash, o.baseHash) {
        return nil  // No recompilation needed
    }
    // Invalidate and recompile only when config changes
}
```

**Key Pattern — Regeneration Levels:**

Endpoint regeneration uses levels to minimize work:

```go
datapathRegenCtxt.bpfHeaderfilesHash = e.orchestrator.EndpointHash(e)
if datapathRegenCtxt.bpfHeaderfilesHash != e.bpfHeaderfileHash {
    datapathRegenCtxt.regenerationLevel = regeneration.RegenerateWithDatapath
}
// Levels: no-op < config-only < full recompile
```

**Why it's clever:** Content-addressed caching means only the *first* endpoint with
a given configuration pays the compilation cost. Subsequent endpoints with the same
config get the cached BPF object for free.

---

## 14. Set Container: Single-Element Optimization

**Location:** `pkg/container/set/set.go`

**Problem:** Many sets in the control plane contain exactly one element. Allocating
a `map` for each is wasteful.

**Solution:** A generic set that uses a pointer for single elements and lazily
allocates a map only when needed:

```go
type Set[T comparable] struct {
    single  *T            // Optimize: 1-element set uses pointer
    members map[T]empty   // Multi-element set uses map
}

func (s *Set[T]) Insert(member T) bool {
    switch s.Len() {
    case 0:
        s.single = &member  // No map allocation
        return true
    case 1:
        s.members = make(map[T]empty, 2)  // Lazy allocation
        s.members[*s.single] = empty{}
        s.single = nil
    }
    s.members[member] = empty{}
    return true
}
```

**Why it's clever:** This optimization is invisible to callers but significantly
reduces memory allocation in the common case. The pattern generalizes well to
any collection type where the single-element case dominates.

---

## 15. StoppableWaitGroup: Graceful Shutdown Primitive

**Location:** `pkg/lock/stoppable_waitgroup.go`

**Problem:** Standard `sync.WaitGroup` panics if `Add()` is called after `Wait()`.
During graceful shutdown, new work may still be submitted while shutdown is in
progress.

**Solution:** A wait group that transitions to no-op mode after `Stop()`:

```go
func (l *StoppableWaitGroup) Add() DoneFunc {
    select {
    case <-l.noopAdd:       // Already stopped?
        return func() {}    // Return no-op
    default:
        l.i.Add(1)
        var once sync.Once
        return func() {
            once.Do(l.done)  // Can't call done() twice
        }
    }
}

func (l *StoppableWaitGroup) Stop() {
    l.stopOnce.Do(func() {
        done := l.Add()   // Add one for the close itself
        close(l.noopAdd)  // Future Add() calls are no-ops
        done()            // Trigger cleanup if counter ≤ 0
    })
}
```

**Why it's clever:** Returns a `DoneFunc` instead of requiring callers to remember
to call `Done()` — each done is wrapped in `sync.Once` preventing double-free bugs.
The channel-based state machine elegantly handles the race between `Add()` and
`Stop()`.

---

## 16. OnDemand: Reference-Counted Lazy Resource Lifecycle

**Location:** `pkg/hive/ondemand.go`

**Problem:** Expensive resources (database connections, gRPC clients) should only
be started when actually needed and stopped when no longer used.

**Solution:** A generic wrapper that starts a resource on first `Acquire()` and
stops it on last `Release()`:

```go
type OnDemand[Resource any] interface {
    Acquire(context.Context) (Resource, error)
    Release(Resource) error
}

func (o *onDemand[Resource]) Acquire(ctx context.Context) (Resource, error) {
    o.mu.Lock()
    if o.refCount == 0 {
        if err := o.lc.Start(o.log, ctx); err != nil {
            return r, fmt.Errorf("failed to start: %w", err)
        }
    }
    o.refCount++
    o.mu.Unlock()
    return o.resource, nil
}

func (o *onDemand[Resource]) Release(r Resource) error {
    o.mu.Lock()
    o.refCount--
    if o.refCount == 0 {
        o.lc.Stop(o.log, context.Background())
    }
    o.mu.Unlock()
    return nil
}
```

**Why it's clever:** Integrates directly with Hive's lifecycle hooks, so the
resource's start/stop hooks are automatically managed. The generic type parameter
ensures type safety without interface casts.

---

## Cross-Cutting Themes

These patterns share several recurring principles:

| Principle | Examples |
|-----------|----------|
| **Minimize critical sections** | StateDB two-phase commit, SelectorCache async notifications, Policy lock minimization |
| **Content-addressed caching** | BPF object cache, endpoint regeneration levels |
| **Lazy initialization** | Resource promise pattern, OnDemand, lazy notification handler |
| **Atomic pointer swap for lock-free reads** | StateDB root swap, SelectorCache `readTxn`, Policy `cachedSelectorPolicy` |
| **Channel-based coordination** | Watch channels, trigger wakeup, StoppableWaitGroup |
| **Fail-fast on misuse** | Fence panic on late Add(), Resource finalizer on forgotten Done(), Module ID regex |
| **Cluster-aware scaling** | Backoff logarithmic scaling, etcd lease pooling |
| **Reference counting** | etcd lease manager, OnDemand lifecycle, StoppableWaitGroup |
| **Event folding** | Trigger mechanism, workqueue deduplication, SelectorCache batched commits |

---

## Recommended Reading Order

For studying these patterns, the recommended order is:

1. **Promise** (`pkg/promise/`) — simplest, foundational async primitive
2. **Trigger** (`pkg/trigger/`) — event folding, widely used
3. **Hive** (`pkg/hive/`) — the DI framework everything builds on
4. **StateDB** (`vendor/github.com/cilium/statedb/`) — the in-process database
5. **Reconciler** (`vendor/github.com/cilium/statedb/reconciler/`) — desired-state engine on top of StateDB
6. **Resource** (`pkg/k8s/resource/`) — Kubernetes integration layer
7. **SelectorCache** (`pkg/policy/selectorcache.go`) — policy engine internals
8. **Lock patterns** (`pkg/lock/`) — concurrency primitives

Each layer builds on the previous ones, making the learning curve manageable.
