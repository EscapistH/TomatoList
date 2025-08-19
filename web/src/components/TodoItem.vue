<template>
    <div class="task-item" :class="{ completed: task.completed }">
        <input type="checkbox" v-model="task.completed" class="task-checkbox" />
        <div class="task-content">
            <span class="task-text">{{ task.text }}</span>
            <span class="task-date">{{ task.createdAt.toLocaleString() }}</span>
        </div>
        <button @click="handleRemove" class="remove-button">删除</button>
    </div>
</template>

<script setup lang="ts">
interface Task {
    id: number;
    text: string;
    completed: boolean;
    createdAt: Date;
}

const props = defineProps<{
    task: Task;
}>();

const emit = defineEmits<{
    (e: "remove", id: number): void;
}>();

const handleRemove = () => {
    emit("remove", props.task.id);
};
</script>

<style scoped>
.task-item {
    display: flex;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid #eee;
    gap: 1rem;
}

.task-item:last-child {
    border-bottom: none;
}

.task-item.completed .task-text {
    text-decoration: line-through;
    color: #888;
}

.task-checkbox {
    transform: scale(1.2);
}

.task-content {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.task-text {
    font-size: 1.1rem;
    margin-bottom: 0.2rem;
}

.task-date {
    color: #888;
    font-size: 0.8rem;
}

.remove-button {
    padding: 0.5rem 1rem;
    background-color: #6c757d;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    height: fit-content;
    align-self: flex-start;
}

.remove-button:hover {
    background-color: #5a6268;
}
</style>
