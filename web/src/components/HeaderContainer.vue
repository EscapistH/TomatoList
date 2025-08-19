<template>
    <header class="header-container">
        <div class="header-content">
            <h1 class="title">TomatoList</h1>
            <nav class="navigation">
                <ul>
                    <li><a href="#" class="nav-link">首页</a></li>
                    <li><a href="#" class="nav-link">任务</a></li>
                    <li><a href="#" class="nav-link">设置</a></li>
                </ul>
            </nav>
        </div>
    </header>
</template>

<script setup lang="ts">
// Header组件不需要复杂的逻辑，使用基础设置即可
</script>

<style scoped>
.header-container {
    background-color: #ff6b6b;
    color: white;
    padding: 1rem 2rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
}

.title {
    margin: 0;
    font-size: 1.8rem;
    font-weight: bold;
}

.navigation ul {
    display: flex;
    list-style: none;
}

.navigation li {
    margin-left: 1.5rem;
}

.nav-link {
    color: white;
    text-decoration: none;
    font-weight: 500;
    transition: opacity 0.3s;
}

.nav-link:hover {
    opacity: 0.8;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .header-content {
        flex-direction: column;
    }

    .title {
        margin-bottom: 1rem;
    }

    .navigation ul {
        padding: 0;
    }

    .navigation li {
        margin: 0 0.5rem;
    }
}
</style>
