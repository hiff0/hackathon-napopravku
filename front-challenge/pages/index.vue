<template>
    <div class="home-container">
        <h1>Присоединившихся пользователей: {{ usersCount }}</h1>
        <div v-if="usersCount >= 2" class="ready-button">
        <NButton @click="onReady">Готов</NButton>
        <p v-if="isReady">Ожидание второго игрока...</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const usersCount = ref(0);
const isReady = ref(false);
let socket : WebSocket | null = null;

onMounted(() => {

    socket  = new WebSocket('ws://localhost:8081/ws');

    socket.onopen = () => {
        console.info('Websocket opened')
    }

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        usersCount.value = data.usersCount;
        console.log('message: ', data);
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed');
    };
});

onBeforeUnmount(() => {
    if (socket) { 
        socket.close()
    };
});

const onReady = () => {
    isReady.value = true;
    checkIfBothReady();
};

const checkIfBothReady = () => {
    setTimeout(() => {
        if (usersCount.value >= 2) {
        const randomId = Math.random().toString(36).substring(7);
        router.push(`/challenge/${randomId}`);
        }
    }, 2000);
};
</script>

<style scoped>
.home-container {
    text-align: center;
    padding: 20px;
}

.ready-button {
    margin-top: 20px;
}

.ui-button {
    font-size: 1.2rem;
    padding: 10px 20px;
}
</style>
