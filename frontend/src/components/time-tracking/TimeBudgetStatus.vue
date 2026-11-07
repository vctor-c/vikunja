<template>
	<div
		v-if="estimate > 0 && summary"
		v-cy="'timeBudgetStatus'"
		class="time-budget-status"
		:class="`has-text-${status}`"
		:data-status="status"
		role="status"
	>
		<div>{{ $t('timeTracking.budget.total', {spent: hours(spent), estimated: hours(estimate)}) }}</div>
		<div>{{ $t(spent > estimate ? 'timeTracking.budget.exceeded' : 'timeTracking.budget.remaining', {hours: hours(Math.abs(estimate - spent))}) }}</div>
	</div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {useNow} from '@vueuse/core'
import {useI18n} from 'vue-i18n'
import {serverNowSeconds} from '@/helpers/serverClock'
import type {Task} from '@/client/generated'

const props = defineProps<{task?: Task}>()
const {locale} = useI18n()
const now = useNow({interval: 1000})
const estimate = computed(() => props.task?.time_estimate_seconds ?? 0)
const summary = computed(() => props.task?.time_tracking)
const spent = computed(() => {
	// Depend on the shared clock tick, using the authenticated server clock offset.
	now.value.getTime()
	const value = summary.value
	if (!value) return 0
	const elapsed = Math.max(0, serverNowSeconds() - new Date(value.calculated_at ?? '').getTime() / 1000)
	return (value.tracked_seconds ?? 0) + elapsed * (value.running_timers ?? 0)
})
const status = computed(() => spent.value > estimate.value ? 'danger' : spent.value >= estimate.value * 0.8 ? 'warning' : 'success')
function hours(seconds: number) {
	return (seconds / 3600).toLocaleString(locale.value, {maximumFractionDigits: 2})
}
</script>
