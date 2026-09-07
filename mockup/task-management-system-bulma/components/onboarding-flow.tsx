'use client'

import { useEffect, useMemo, useState, type KeyboardEvent } from 'react'
import {
  ArrowLeft,
  ArrowRight,
  Check,
  CheckCircle2,
  Database,
  Folder,
  Globe,
  HardDrive,
  Inbox,
  Layers,
  ListChecks,
  Loader2,
  Puzzle,
  Rocket,
  Server,
  Sparkles,
  User as UserIcon,
} from 'lucide-react'

export type HostingMode = 'local' | 'self-hosted' | 'standalone-web' | 'embedded'

export interface OnboardingData {
  hostingMode: HostingMode
  persistenceLocation: string
  spaceName: string
  userName: string
  firstListName: string
}

type StepId =
  | 'welcome'
  | 'hosting-mode'
  | 'persistence-location'
  | 'initialize'
  | 'space'
  | 'user'
  | 'system-setup'
  | 'first-list'
  | 'complete'

type Phase = 'instance' | 'workspace' | 'ready'

const PHASE_OF: Record<StepId, Phase> = {
  welcome: 'instance',
  'hosting-mode': 'instance',
  'persistence-location': 'instance',
  initialize: 'instance',
  space: 'workspace',
  user: 'workspace',
  'system-setup': 'workspace',
  'first-list': 'workspace',
  complete: 'ready',
}

const PHASES: { id: Phase; label: string }[] = [
  { id: 'instance', label: 'Instance' },
  { id: 'workspace', label: 'Workspace' },
  { id: 'ready', label: 'Ready' },
]

const HOSTING_OPTIONS: {
  value: HostingMode
  title: string
  desc: string
  persistence: string
  icon: typeof Database
}[] = [
  {
    value: 'local',
    title: 'Local',
    desc: 'Runs directly on this machine, with local filesystem-backed persistence.',
    persistence: 'SQLite',
    icon: HardDrive,
  },
  {
    value: 'self-hosted',
    title: 'Self-hosted',
    desc: 'Runs on infrastructure you control — Docker, a NAS, home server, or VPS.',
    persistence: 'Server database',
    icon: Server,
  },
  {
    value: 'standalone-web',
    title: 'Standalone Web',
    desc: 'Runs independently in the browser as a PWA or standalone app.',
    persistence: 'IndexedDB / OPFS',
    icon: Globe,
  },
  {
    value: 'embedded',
    title: 'Embedded',
    desc: 'Runs inside another application as a component or SDK, using the host for services.',
    persistence: 'Host-provided',
    icon: Puzzle,
  },
]

// Persistence type is dictated by the hosting mode.
function persistenceForMode(mode: HostingMode): string {
  return HOSTING_OPTIONS.find((o) => o.value === mode)?.persistence ?? 'SQLite'
}

function hostingLabel(mode: HostingMode): string {
  return HOSTING_OPTIONS.find((o) => o.value === mode)?.title ?? 'Local'
}

// A data directory is only needed for local and self-hosted modes.
function needsLocation(mode: HostingMode): boolean {
  return mode === 'local' || mode === 'self-hosted'
}

function submitOnEnter(e: KeyboardEvent, fn: () => void) {
  // Respect IME composition (CJK) and Safari's unreliable final event.
  if (e.key !== 'Enter') return
  if (e.nativeEvent.isComposing || e.keyCode === 229) return
  e.preventDefault()
  fn()
}

export function OnboardingFlow({ onComplete }: { onComplete: (data: OnboardingData) => void }) {
  const [step, setStep] = useState<StepId>('welcome')
  const [autoReady, setAutoReady] = useState(false)
  const [data, setData] = useState<OnboardingData>({
    hostingMode: 'local',
    persistenceLocation: '~/listello',
    spaceName: 'Personal',
    userName: '',
    firstListName: 'Errands',
  })

  // The location step only applies to local and self-hosted modes (per the board's note).
  const order = useMemo<StepId[]>(() => {
    const all: StepId[] = [
      'welcome',
      'hosting-mode',
      'persistence-location',
      'initialize',
      'space',
      'user',
      'system-setup',
      'first-list',
      'complete',
    ]
    return all.filter((id) => id !== 'persistence-location' || needsLocation(data.hostingMode))
  }, [data.hostingMode])

  const index = order.indexOf(step)
  const isAutoStep = step === 'initialize' || step === 'system-setup'

  // Reset the "auto complete" gate whenever we arrive on an automated step.
  useEffect(() => {
    setAutoReady(false)
  }, [step])

  function goNext() {
    const next = order[index + 1]
    if (next) setStep(next)
    else onComplete(data)
  }

  function goBack() {
    const prev = order[index - 1]
    if (prev) setStep(prev)
  }

  // Proceed without creating a first list: clear the name so no list is created.
  function skipList() {
    setData((d) => ({ ...d, firstListName: '' }))
    goNext()
  }

  const canAdvance = (() => {
    if (isAutoStep) return autoReady
    if (step === 'persistence-location') return data.persistenceLocation.trim().length > 0
    if (step === 'space') return data.spaceName.trim().length > 0
    if (step === 'user') return data.userName.trim().length > 0
    if (step === 'first-list') return data.firstListName.trim().length > 0
    return true
  })()

  const activePhase = PHASE_OF[step]

  function segFill(phase: Phase): number {
    const phaseSteps = order.filter((id) => PHASE_OF[id] === phase)
    if (phaseSteps.length === 0) return 0
    const passed = phaseSteps.filter((id) => order.indexOf(id) < index).length
    const onIt = phaseSteps.includes(step)
    const done = passed + (onIt ? 1 : 0)
    // Fully-passed phases read as complete.
    const allPassed = order.indexOf(phaseSteps[phaseSteps.length - 1]) < index
    if (allPassed) return 100
    return Math.min(100, Math.round((done / phaseSteps.length) * 100))
  }

  return (
    <div className="onboarding-page">
      <header className="onboarding-topbar">
        <div className="onboarding-rail">
          <div className="brand-row">
            <span className="brand-mark">
              <span className="brand-dot">
                <Check size={18} strokeWidth={3} />
              </span>
              Listello
            </span>
            <span className="setup-badge">First-time setup</span>
          </div>

          <div className="phase-progress" aria-hidden>
            {PHASES.map((p, i) => {
              const fill = segFill(p.id)
              const isActive = activePhase === p.id
              const isDone = fill === 100 && !isActive
              const state = isDone ? 'is-done' : isActive ? 'is-active' : 'is-upcoming'
              return (
                <div key={p.id} className={`phase-seg ${state}`}>
                  <span className="phase-seg-head">
                    <span className="phase-seg-index">
                      {isDone ? <Check size={12} strokeWidth={3} /> : i + 1}
                    </span>
                    <span className="phase-seg-label">{p.label}</span>
                  </span>
                  <span className="phase-seg-track">
                    <span className="phase-seg-fill" style={{ width: `${fill}%` }} />
                  </span>
                </div>
              )
            })}
          </div>
        </div>
      </header>

      <div className="onboarding-body">
        <div className="onboarding-inner">
          {step === 'welcome' && (
            <div>
              <span className="big-icon">
                <Sparkles size={26} />
              </span>
              <h1 className="step-title text-balance">Welcome to Listello</h1>
              <p className="step-lead text-pretty">
                {"Let's create your instance and set up a calm, GTD-style workspace. "}
                {"It only takes a minute, and you can change everything later."}
              </p>
              <div className="step-content">
                <div className="setup-check is-done">
                  <span className="setup-check-status">
                    <Server size={18} />
                  </span>
                  <span className="setup-check-label">Choose how Listello is hosted</span>
                </div>
                <div className="setup-check is-done">
                  <span className="setup-check-status">
                    <Layers size={18} />
                  </span>
                  <span className="setup-check-label">Create your space and profile</span>
                </div>
                <div className="setup-check is-done">
                  <span className="setup-check-status">
                    <ListChecks size={18} />
                  </span>
                  <span className="setup-check-label">Start your first list</span>
                </div>
              </div>
            </div>
          )}

          {step === 'hosting-mode' && (
            <div>
              <p className="step-eyebrow">Instance · Hosting</p>
              <h1 className="step-title text-balance">How should Listello be hosted?</h1>
              <p className="step-lead text-pretty">
                Your hosting mode decides where this instance runs and how your data is stored. You
                can migrate later.
              </p>
              <div className="step-content choice-list">
                {HOSTING_OPTIONS.map((opt) => {
                  const Icon = opt.icon
                  const selected = data.hostingMode === opt.value
                  return (
                    <button
                      key={opt.value}
                      type="button"
                      className={`choice-card ${selected ? 'is-selected' : ''}`}
                      aria-pressed={selected}
                      onClick={() => setData((d) => ({ ...d, hostingMode: opt.value }))}
                    >
                      <span className="choice-icon">
                        <Icon size={20} />
                      </span>
                      <span>
                        <span className="choice-title is-block">{opt.title}</span>
                        <span className="choice-desc">{opt.desc}</span>
                        <span className="choice-meta">
                          <Database size={12} />
                          Uses {opt.persistence}
                        </span>
                      </span>
                      {selected && (
                        <span className="choice-check">
                          <CheckCircle2 size={20} />
                        </span>
                      )}
                    </button>
                  )
                })}
              </div>
            </div>
          )}

          {step === 'persistence-location' && (
            <div>
              <p className="step-eyebrow">Instance · {hostingLabel(data.hostingMode)}</p>
              <h1 className="step-title text-balance">Choose a data directory</h1>
              <p className="step-lead text-pretty">
                {"This is where Listello keeps its data for this instance. You can move it later."}
              </p>
              <div className="step-content">
                <div className="field">
                  <label className="label" htmlFor="onb-location">
                    Data directory
                  </label>
                  <div className="control has-icons-left">
                    <input
                      id="onb-location"
                      className="input"
                      type="text"
                      value={data.persistenceLocation}
                      placeholder="~/listello"
                      onChange={(e) =>
                        setData((d) => ({ ...d, persistenceLocation: e.target.value }))
                      }
                      onKeyDown={(e) => submitOnEnter(e, () => canAdvance && goNext())}
                    />
                    <span className="icon is-small is-left">
                      <Folder size={16} />
                    </span>
                  </div>
                  <p className="help">
                    {`Logs, the ${persistenceForMode(data.hostingMode)} database, and config are stored here.`}
                  </p>
                </div>
              </div>
            </div>
          )}

          {step === 'initialize' && (
            <div>
              <p className="step-eyebrow">Instance · Initialize</p>
              <h1 className="step-title text-balance">Setting up persistence</h1>
              <p className="step-lead text-pretty">
                {"We're initializing the "}
                <strong>{persistenceForMode(data.hostingMode)}</strong>
                {" store for your "}
                <strong>{hostingLabel(data.hostingMode).toLowerCase()}</strong>
                {" instance."}
              </p>
              <div className="step-content">
                <AutoSetup
                  key="initialize"
                  steps={[
                    `Create instance`,
                    `Prepare ${persistenceForMode(data.hostingMode)} store`,
                    `Initialize persistence`,
                  ]}
                  onDone={() => setAutoReady(true)}
                />
              </div>
            </div>
          )}

          {step === 'space' && (
            <div>
              <p className="step-eyebrow">Workspace · Space</p>
              <h1 className="step-title text-balance">Name your space</h1>
              <p className="step-lead text-pretty">
                A space groups your lists together. Most people start with a single personal space.
              </p>
              <div className="step-content">
                <div className="field">
                  <label className="label" htmlFor="onb-space">
                    Space name
                  </label>
                  <div className="control has-icons-left">
                    <input
                      id="onb-space"
                      className="input"
                      type="text"
                      value={data.spaceName}
                      placeholder="Personal"
                      onChange={(e) => setData((d) => ({ ...d, spaceName: e.target.value }))}
                      onKeyDown={(e) => submitOnEnter(e, () => canAdvance && goNext())}
                    />
                    <span className="icon is-small is-left">
                      <Layers size={16} />
                    </span>
                  </div>
                </div>
              </div>
            </div>
          )}

          {step === 'user' && (
            <div>
              <p className="step-eyebrow">Workspace · You</p>
              <h1 className="step-title text-balance">{"What should we call you?"}</h1>
              <p className="step-lead text-pretty">
                {"Your name shows up on comments and activity. It's just for you — no account needed."}
              </p>
              <div className="step-content">
                <div className="field">
                  <label className="label" htmlFor="onb-user">
                    Your name
                  </label>
                  <div className="control has-icons-left">
                    <input
                      id="onb-user"
                      className="input"
                      type="text"
                      value={data.userName}
                      placeholder="e.g. Alex"
                      autoFocus
                      onChange={(e) => setData((d) => ({ ...d, userName: e.target.value }))}
                      onKeyDown={(e) => submitOnEnter(e, () => canAdvance && goNext())}
                    />
                    <span className="icon is-small is-left">
                      <UserIcon size={16} />
                    </span>
                  </div>
                </div>
              </div>
            </div>
          )}

          {step === 'system-setup' && (
            <div>
              <p className="step-eyebrow">Workspace · Automatic</p>
              <h1 className="step-title text-balance">Getting things ready</h1>
              <p className="step-lead text-pretty">
                {"Listello is wiring up the essentials for "}
                <strong>{data.spaceName.trim() || 'your space'}</strong>
                {"."}
              </p>
              <div className="step-content">
                <AutoSetup
                  key="system-setup"
                  steps={['Create Inbox', `Assign ${data.spaceName.trim() || 'space'} to ${data.userName.trim() || 'you'}`]}
                  onDone={() => setAutoReady(true)}
                />
              </div>
            </div>
          )}

          {step === 'first-list' && (
            <div>
              <p className="step-eyebrow">Workspace · First list</p>
              <h1 className="step-title text-balance">Create your first list</h1>
              <p className="step-lead text-pretty">
                Lists hold the tasks you want to act on. Give your first one a name, or skip and
                create one later.
              </p>
              <div className="step-content">
                <div className="field">
                  <label className="label" htmlFor="onb-list">
                    List name
                  </label>
                  <div className="control has-icons-left">
                    <input
                      id="onb-list"
                      className="input"
                      type="text"
                      value={data.firstListName}
                      placeholder="Errands"
                      onChange={(e) => setData((d) => ({ ...d, firstListName: e.target.value }))}
                      onKeyDown={(e) => submitOnEnter(e, () => canAdvance && goNext())}
                    />
                    <span className="icon is-small is-left">
                      <ListChecks size={16} />
                    </span>
                  </div>
                  <p className="suggest-hint">Or start with one of these</p>
                  <div className="tags mt-2">
                    {['Errands', 'Shopping', 'Ideas', 'Reading', 'Goals'].map((s) => (
                      <button
                        key={s}
                        type="button"
                        className={`tag is-medium ${data.firstListName === s ? 'is-primary' : ''}`}
                        onClick={() => setData((d) => ({ ...d, firstListName: s }))}
                      >
                        {s}
                      </button>
                    ))}
                  </div>
                </div>
                <button type="button" className="skip-link" onClick={skipList}>
                  Skip for now
                </button>
              </div>
            </div>
          )}

          {step === 'complete' && (
            <div>
              <span className="big-icon">
                <Rocket size={26} />
              </span>
              <h1 className="step-title text-balance">
                {"You're all set, "}
                {data.userName.trim() || 'friend'}
              </h1>
              <p className="step-lead text-pretty">
                Your instance is ready. Here is what we set up — you can change any of it later.
              </p>
              <div className="step-content box">
                <div className="summary-row">
                  <span className="summary-key">Hosting</span>
                  <span className="summary-val">{hostingLabel(data.hostingMode)}</span>
                </div>
                <div className="summary-row">
                  <span className="summary-key">Persistence</span>
                  <span className="summary-val">{persistenceForMode(data.hostingMode)}</span>
                </div>
                {needsLocation(data.hostingMode) && (
                  <div className="summary-row">
                    <span className="summary-key">Location</span>
                    <span className="summary-val is-family-code">{data.persistenceLocation}</span>
                  </div>
                )}
                <div className="summary-row">
                  <span className="summary-key">Space</span>
                  <span className="summary-val">{data.spaceName}</span>
                </div>
                <div className="summary-row">
                  <span className="summary-key">First list</span>
                  <span className={`summary-val ${data.firstListName.trim() ? '' : 'is-muted'}`}>
                    {data.firstListName.trim() || 'Skipped'}
                  </span>
                </div>
                <div className="summary-row">
                  <span className="summary-key">Inbox</span>
                  <span className="summary-val">
                    <span className="icon-text">
                      <span className="icon has-text-primary">
                        <Inbox size={16} />
                      </span>
                      <span>Ready</span>
                    </span>
                  </span>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      <footer className="onboarding-footer">
        {index > 0 && step !== 'complete' && (
          <button type="button" className="button is-light" onClick={goBack}>
            <span className="icon">
              <ArrowLeft size={18} />
            </span>
            <span>Back</span>
          </button>
        )}
        <button
          type="button"
          className="button is-primary footer-grow"
          disabled={!canAdvance}
          onClick={goNext}
        >
          <span>{footerLabel(step)}</span>
          <span className="icon">
            {step === 'complete' ? <Rocket size={18} /> : <ArrowRight size={18} />}
          </span>
        </button>
      </footer>
    </div>
  )
}

function footerLabel(step: StepId): string {
  if (step === 'welcome') return 'Create instance'
  if (step === 'initialize') return 'Continue'
  if (step === 'first-list') return 'Create list'
  if (step === 'complete') return 'Enter Listello'
  return 'Continue'
}

/**
 * Plays through a short checklist to represent an automated, System-driven
 * step (Initialize persistence, Create Inbox, Assign space). Calls onDone
 * once every line has "completed".
 */
function AutoSetup({ steps, onDone }: { steps: string[]; onDone: () => void }) {
  const [doneCount, setDoneCount] = useState(0)

  useEffect(() => {
    setDoneCount(0)
    let current = 0
    const timers: ReturnType<typeof setTimeout>[] = []
    const tick = () => {
      current += 1
      setDoneCount(current)
      if (current >= steps.length) {
        onDone()
        return
      }
      timers.push(setTimeout(tick, 650))
    }
    timers.push(setTimeout(tick, 650))
    return () => timers.forEach(clearTimeout)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [steps.length])

  return (
    <div>
      {steps.map((label, i) => {
        const done = i < doneCount
        const active = i === doneCount
        return (
          <div
            key={label}
            className={`setup-check ${done ? 'is-done' : active ? '' : 'is-pending'}`}
          >
            <span className="setup-check-status">
              {done ? (
                <CheckCircle2 size={18} />
              ) : active ? (
                <Loader2 size={18} className="spin" />
              ) : (
                <span className="check-toggle" style={{ width: 16, height: 16 }} />
              )}
            </span>
            <span className="setup-check-label">{label}</span>
          </div>
        )
      })}
    </div>
  )
}
