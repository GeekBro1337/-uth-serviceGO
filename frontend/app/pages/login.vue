<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900 px-4 py-8"
  >
    <div
      class="w-full max-w-md sm:max-w-lg lg:max-w-xl min-w-[340px] bg-white dark:bg-gray-800 rounded-2xl shadow-2xl overflow-hidden"
    >
      <!-- Header -->
      <div class="p-8 border-b border-gray-200 dark:border-gray-700 w-full">
        <div class="text-center">
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Добро пожаловать
          </h1>
          <p class="text-gray-600 dark:text-gray-400">Войдите в свой аккаунт</p>
        </div>
      </div>

      <!-- Content -->
      <div class="p-8">
        <UAlert
          v-if="error"
          color="error"
          variant="soft"
          :title="error"
          class="mb-4"
          icon="i-heroicons-exclamation-triangle"
        />

        <form @submit.prevent="login" class="space-y-4">
          <UFormGroup label="Имя пользователя" name="username" required>
            <UInput
              v-model="username"
              type="text"
              placeholder="Введите имя пользователя"
              size="xl"
              :disabled="loading"
              icon="i-heroicons-user"
              autocomplete="username"
            />
          </UFormGroup>

          <UFormGroup label="Пароль" name="password" required>
            <UInput
              v-model="password"
              type="password"
              placeholder="Введите пароль"
              size="xl"
              :disabled="loading"
              icon="i-heroicons-lock-closed"
              autocomplete="current-password"
            />
          </UFormGroup>

          <UButton
            type="submit"
            block
            size="lg"
            :loading="loading"
            :disabled="!username || !password"
            class="mt-6"
          >
            <template v-if="!loading">Войти</template>
            <template v-else>Вход...</template>
          </UButton>
        </form>
      </div>

      <!-- Footer -->
      <div class="p-8 border-t border-gray-200 dark:border-gray-700">
        <div class="text-center text-sm text-gray-600 dark:text-gray-400">
          Нет аккаунта?
          <NuxtLink
            to="/register"
            class="font-semibold text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 transition-colors ml-1"
          >
            Зарегистрироваться
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { accessToken, refreshToken } = useApi();
const router = useRouter();
const config = useRuntimeConfig();
const baseURL = config.public.apiBase || "http://localhost:8999";

const username = ref("");
const password = ref("");
const loading = ref(false);
const error = ref<string | null>(null);

const login = async () => {
  if (!username.value || !password.value) {
    error.value = "Заполните все поля";
    return;
  }

  loading.value = true;
  error.value = null;

  try {
    const res: any = await $fetch(baseURL + "/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include", // важно для получения cookies
      body: JSON.stringify({
        username: username.value,
        password: password.value,
      }),
    });

    // В новой версии ответ содержит только "token" (access token)
    // Refresh токен автоматически устанавливается в HttpOnly cookie на сервере
    if (!res.token) {
      error.value = "Ошибка: токен не получен";
      return;
    }

    accessToken.value = res.token;

    await new Promise((resolve) => setTimeout(resolve, 100));
    await router.push("/");
  } catch (e: any) {
    console.error("Ошибка логина:", e);
    error.value =
      e.data?.error || e.message || "Ошибка входа. Проверьте данные.";
  } finally {
    loading.value = false;
  }
};
</script>
