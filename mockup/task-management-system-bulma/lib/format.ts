import type { Priority } from './types'

export const PRIORITIES: { value: Priority; label: string; dot: string }[] = [
  { value: 'none', label: 'No priority', dot: 'hsl(0deg 0% 60%)' },
  { value: 'low', label: 'Low', dot: 'hsl(0deg 0% 45%)' },
  { value: 'medium', label: 'Medium', dot: 'hsl(38deg 92% 50%)' },
  { value: 'high', label: 'High', dot: 'hsl(348deg 86% 55%)' },
]

export function priorityMeta(p: Priority) {
  return PRIORITIES.find((x) => x.value === p) ?? PRIORITIES[0]
}

export function formatDue(iso: string | null): { label: string; tone: 'past' | 'today' | 'soon' | 'none' } {
  if (!iso) return { label: '', tone: 'none' }
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const due = new Date(iso + 'T00:00:00')
  const diff = Math.round((due.getTime() - today.getTime()) / 86400000)
  if (diff < 0) return { label: diff === -1 ? 'Yesterday' : `${Math.abs(diff)}d overdue`, tone: 'past' }
  if (diff === 0) return { label: 'Today', tone: 'today' }
  if (diff === 1) return { label: 'Tomorrow', tone: 'soon' }
  if (diff < 7) return { label: `${diff}d`, tone: 'soon' }
  return {
    label: due.toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    tone: 'none',
  }
}

export function relativeTime(ts: number): string {
  const diff = Date.now() - ts
  const mins = Math.round(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.round(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.round(hours / 24)
  if (days < 7) return `${days}d ago`
  return new Date(ts).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}
