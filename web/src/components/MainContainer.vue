<template>
    <main class="main-container">
        <div class="content-wrapper">
            <section class="task-input-section">
                <input
                    type="text"
                    class="task-input"
                    placeholder="添加新任务..."
                    v-model="newTask"
                    @keyup.enter="addTask"
                />
                <button class="add-button" @click="addTask">添加</button>
            </section>

            <section class="task-list-section">
                <h2>待办事项</h2>

                <div class="task-list" v-show="tasks.length > 0">
                    <TodoItem
                        v-for="task in tasks"
                        :key="task.id"
                        :task="task"
                        @remove="removeTask"
                    />
                </div>

                <!-- 当有任务时显示统计信息 -->
                <div class="task-stats" v-if="tasks.length > 0">
                    <p>
                        已完成:{{ completedTasksCount }} | 总计:
                        {{ tasks.length }}
                    </p>
                </div>

                <!-- 当没有任务时显示提示信息 -->
                <div class="empty-state" v-else>
                    <p>暂无任务，添加一个新任务开始吧！</p>
                </div>
            </section>
        </div>
    </main>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import TodoItem from "./TodoItem.vue";

interface Task {
    id: number;
    text: string;
    completed: boolean;
    createdAt: Date;
}

const newTask = ref("");
const tasks = ref<Task[]>([]);

const addTask = () => {
    if (newTask.value.trim() !== "") {
        tasks.value.push({
            id: Date.now(),
            text: newTask.value.trim(),
            completed: false,
            createdAt: new Date(),
        });
        newTask.value = "";
    }
};

const removeTask = (id: number) => {
    tasks.value = tasks.value.filter((task) => task.id !== id);
};

const completedTasksCount = computed(() => {
    return tasks.value.filter((task) => task.completed).length;
});
</script>

<style scoped>
.main-container {
    padding: 2rem;
    background-color: #f8f9fa;
    min-height: calc(100vh - 80px);
}

.content-wrapper {
    max-width: 1200px;
    margin: 0 auto;
}

.task-input-section {
    display: flex;
    margin-bottom: 2rem;
    gap: 1rem;
}

.task-input {
    flex: 1;
    padding: 0.75rem;
    border: 2px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
}

.task-input:focus {
    outline: none;
    border-color: #ff6b6b;
}

.add-button {
    padding: 0.75rem 1.5rem;
    background-color: #ff6b6b;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    font-weight: bold;
    transition: background-color 0.3s;
}

.add-button:hover {
    background-color: #ff5252;
}

.task-list-section h2 {
    color: #333;
    margin-bottom: 1rem;
}

.task-list {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 1rem;
}

.task-stats {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid #eee;
    text-align: center;
    color: #6c757d;
}

.empty-state {
    text-align: center;
    color: #6c757d;
    padding: 2rem;
}

.empty-state p {
    margin: 0;
    font-style: italic;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .main-container {
        padding: 1rem;
    }

    .task-input-section {
        flex-direction: column;
    }
}
</style>
