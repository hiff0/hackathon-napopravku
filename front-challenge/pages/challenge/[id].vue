<template>
    <div class="challenge-container">
        <h1>Соревнование #{{ id }}</h1>
        <div class="editor-container">
        <!-- Здесь будет редактор кода -->
        <p>Текстовый редактор</p>
        </div>
        <div class="layout-container">
        <!-- Макет для верстки -->
        <p>Макет</p>
        </div>
        <div class="result-container">
        <!-- Результат работы -->
        <p>Результат</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute();
const id = route.params.id; // Получаем ID соревнования
</script>

<style scoped>
.challenge-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
}

.editor-container, .layout-container, .result-container {
    border: 1px solid #ccc;
    padding: 20px;
    margin: 10px;
    width: 80%;
    height: 300px;
}
</style>
