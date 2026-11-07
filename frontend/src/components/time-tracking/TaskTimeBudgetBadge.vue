<template>
	<TimeBudgetStatus
		v-if="enabled"
		:task="budgetTask"
		compact
	/>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {useQuery} from '@tanstack/vue-query'
import type {Task} from '@/client/generated'
import {taskQuery} from '@/client/queries/tasks'
import {useAuthStore} from '@/stores/auth'
import {useConfigStore} from '@/stores/config'
import {PRO_FEATURE} from '@/constants/proFeatures'
import TimeBudgetStatus from './TimeBudgetStatus.vue'

const props = defineProps<{task: Task}>()
const auth = useAuthStore()
const config = useConfigStore()
const enabled = computed(() => (props.task.time_estimate_seconds ?? 0) > 0
	&& !!props.task.id && !auth.isLinkShareAuth && config.isProFeatureEnabled(PRO_FEATURE.TIME_TRACKING))
const {data} = useQuery(computed(() => ({
	...taskQuery(props.task.id ?? 0, ['time_tracking']),
	enabled: enabled.value,
	refetchInterval: enabled.value ? 30000 : false,
})))
// Keep the card's estimate authoritative when it changes through a task edit.
const budgetTask = computed(() => ({...props.task, time_tracking: data.value?.time_tracking ?? props.task.time_tracking}))
</script>
