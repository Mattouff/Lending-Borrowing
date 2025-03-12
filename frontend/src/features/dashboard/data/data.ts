import {
    IconAi,
    IconCash,
    IconDrone,
    IconWebhook
} from '@tabler/icons-react'
import { UserStatus } from './schema'

export const callTypes = new Map<UserStatus, string>([
    ['active', 'bg-teal-100/30 text-teal-900 dark:text-teal-200 border-teal-200'],
    ['inactive', 'bg-neutral-300/40 border-neutral-300'],
    ['invited', 'bg-sky-200/40 text-sky-900 dark:text-sky-100 border-sky-300'],
    [
        'suspended',
        'bg-destructive/10 dark:bg-destructive/50 text-destructive dark:text-primary border-destructive/10',
    ],
])

export const degreesTypes = [
    {
        label: 'Master - Ingénierie de la Blockchain',
        value: 'blockchain',
        icon: IconCash,
    },
    {
        label: 'Master - Ingénierie en Intelligence Artificielle',
        value: 'ia',
        icon: IconAi,
    },
    {
        label: 'Master - Ingénierie du web',
        value: 'manager',
        icon: IconWebhook,
    },
    {
        label: 'Master - Ingénierie des systèmes embarqués',
        value: 'embarque',
        icon: IconDrone,
    },
] as const

export const degreesStatus = [
    {
        label: 'En cours',
        value: 'en_cours',
    },
    {
        label: 'Validé',
        value: 'valide',
    },
    {
        label: 'Refusé',
        value: 'refuse',
    },
]

export const degreesYears = [
    {
        label: '2021',
        value: '2021',
    },
    {
        label: '2022',
        value: '2022',
    },
    {
        label: '2023',
        value: '2023',
    },
    {
        label: '2024',
        value: '2024',
    },
    {
        label: '2025',
        value: '2025',
    }
]