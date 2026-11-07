<template>
	<div class="task-time-tracking">
		<XButton
			v-if="entries.length > 0"
			v-tooltip="$t('timeTracking.logTime')"
			v-cy="'addTaskTimeEntry'"
			:aria-label="$t('timeTracking.logTime')"
			class="is-pulled-end add-time-entry-button d-print-none"
			:class="{'is-active': showForm}"
			variant="secondary"
			icon="plus"
			:shadow="false"
			@click="showForm = !showForm"
		/>
		<h2 class="task-section-title">
			<span class="icon is-grey">
				<Icon :icon="['far', 'clock']" />
			</span>
			{{ $t('timeTracking.title') }}
		</h2>
		<form
			v-if="canWrite"
			class="mbs-4 mbe-4"
			@submit.prevent="saveEstimate"
		>
			<label :for="`time-estimate-${taskId}`">{{ $t('timeTracking.budget.label') }}</label>
			<input
				:id="`time-estimate-${taskId}`"
				v-model="estimatedHours"
				v-cy="'timeEstimate'"
				class="input"
				type="number"
				min="0"
				step="any"
				:placeholder="$t('timeTracking.budget.optional')"
			>
			<XButton
				v-cy="'saveTimeEstimate'"
				class="mbs-2"
				type="submit"
				:disabled="saveMutation.isPending.value || !task"
			>
				{{ $t('misc.save') }}
			</XButton>
		</form>
		<TimeBudgetStatus
			:task="task"
			class="mbe-4"
		/>
		<TimeEntryForm
			v-if="formVisible"
			:task-id="taskId"
			:entry="editingEntry"
			:recent-entries="entries"
			@saved="onSaved"
			@cancel="editingEntry = null"
		/>
		<TimeEntryList
			class="mbs-4"
			:entries="entries"
			:card="false"
			:empty-text="$t('timeTracking.list.emptyTask')"
			:paged="totalPages > 1"
			hide-label-column
			@edit="editingEntry = $event"
			@delete="onDelete"
		/>
		<PaginationEmit
			v-if="totalPages > 1"
			:total-pages="totalPages"
			:current-page="currentPage"
			@pageChanged="currentPage = $event"
		/>
	</div>
</template>

<script setup lang="ts">
import {ref, computed, watch} from 'vue'
import {useQuery} from '@tanstack/vue-query'
import {taskQuery} from '@/client/queries/tasks'
import {useUpdateTaskMutation} from '@/client/queries/taskMutations'
import TimeBudgetStatus from './TimeBudgetStatus.vue'

import PaginationEmit from '@/components/misc/PaginationEmit.vue'
import TimeEntryForm from '@/components/time-tracking/TimeEntryForm.vue'
import TimeEntryList from '@/components/time-tracking/TimeEntryList.vue'

import {useTimeEntries} from '@/composables/useTimeTracking'
import {useDeleteTimeEntryMutation} from '@/client/queries/timeEntries'

import type {TimeEntryResponse as ITimeEntry} from '@/client/queries/timeEntries'

const props = defineProps<{
	taskId: number
	canWrite?: boolean
}>()

const {data: task} = useQuery(computed(() => ({
	...taskQuery(props.taskId, ['time_tracking']),
	refetchInterval: 30000,
})))
const estimatedHours = ref<string | number>('')
const saveMutation = useUpdateTaskMutation()
watch([() => props.taskId, () => task.value?.time_estimate_seconds], ([, seconds]) => {
	estimatedHours.value = seconds ? seconds / 3600 : ''
}, {immediate: true})
async function saveEstimate() {
	const hours = Number(estimatedHours.value)
	if (!props.canWrite || !task.value || !Number.isFinite(hours) || hours < 0 || saveMutation.isPending.value) return
	try {
		await saveMutation.mutateAsync({id: props.taskId, time_estimate_seconds: Math.round(hours * 3600)})
	} catch {
		// Mutation reports the error and keeps the entered estimate for retry.
	}
}

const currentPage = ref(1)
const {entries, totalPages} = useTimeEntries(() => `task_id = ${props.taskId}`, {page: currentPage})
const deleteMutation = useDeleteTimeEntryMutation()
const editingEntry = ref<ITimeEntry | null>(null)
const showForm = ref(false)

watch(() => props.taskId, () => {
	currentPage.value = 1
})

// Like related tasks: the form is implicit when empty, otherwise behind the +.
const formVisible = computed(() => entries.value.length === 0 || showForm.value || editingEntry.value !== null)

function onSaved() {
	editingEntry.value = null
	showForm.value = false
}

async function onDelete(entry: ITimeEntry) {
	if (deleteMutation.isPending.value) {
		return
	}
	try {
		await deleteMutation.mutateAsync({
			id: entry.id,
			taskId: entry.task_id,
		})
	} catch {
		return
	}
	currentPage.value = Math.min(currentPage.value, Math.max(1, totalPages.value))
}
</script>

<style lang="scss" scoped>
.add-time-entry-button {
	svg {
		transition: transform $transition;
	}

	&.is-active svg {
		transform: rotate(45deg);
	}
}
</style>
