<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-50 via-white to-blue-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900 px-4 py-8"
  >
    <div
      class="w-full max-w-md sm:max-w-lg lg:max-w-xl min-w-[340px] bg-white dark:bg-gray-800 rounded-2xl shadow-2xl overflow-hidden"
    >
      <!-- Header -->
      <div class="p-8 border-b border-gray-200 dark:border-gray-700">
        <div class="text-center">
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Создать аккаунт
          </h1>
          <p class="text-gray-600 dark:text-gray-400">
            Заполните форму для регистрации
          </p>
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

        <UAlert
          v-if="success"
          color="success"
          variant="soft"
          title="Регистрация успешна!"
          description="Вы будете перенаправлены на страницу входа"
          class="mb-4"
          icon="i-heroicons-check-circle"
        />

        <form @submit.prevent="registerUser" class="space-y-4">
          <UFormGroup label="Имя пользователя" name="username" required>
            <UInput
              v-model="username"
              type="text"
              placeholder="Введите имя пользователя"
              size="lg"
              :disabled="loading || success"
              icon="i-heroicons-user"
              autocomplete="username"
            />
          </UFormGroup>

          <UFormGroup label="Пароль" name="password" required>
            <UInput
              v-model="password"
              type="password"
              placeholder="Введите пароль"
              size="lg"
              :disabled="loading || success"
              icon="i-heroicons-lock-closed"
              autocomplete="new-password"
            />
          </UFormGroup>

          <UFormGroup label="Роль (опционально)" name="role">
            <USelect
              v-model="role"
              :options="roleOptions"
              placeholder="Выберите роль"
              size="lg"
              :disabled="loading || success"
              icon="i-heroicons-user-circle"
            />
          </UFormGroup>

          <UButton
            type="submit"
            block
            size="lg"
            :loading="loading"
            :disabled="!username || !password || success"
            class="mt-6"
          >
            <template v-if="!loading && !success">
              Зарегистрироваться
            </template>
            <template v-else-if="loading"> Регистрация... </template>
            <template v-else> Успешно! </template>
          </UButton>
        </form>
      </div>

      <!-- Footer -->
      <div class="p-8 border-t border-gray-200 dark:border-gray-700">
        <div class="text-center text-sm text-gray-600 dark:text-gray-400">
          Уже есть аккаунт?
          <NuxtLink
            to="/login"
            class="font-semibold text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 transition-colors ml-1"
          >
            Войти
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const router = useRouter();
const config = useRuntimeConfig();
const baseURL = config.public.apiBase || "http://localhost:8999";

const username = ref("");
const password = ref("");
const role = ref("");
const loading = ref(false);
const error = ref<string | null>(null);
const success = ref(false);

const roleOptions = [
  { label: "Пользователь", value: "user" },
  { label: "Администратор", value: "admin" },
  { label: "Аудитор", value: "auditor" },
];

const registerUser = async () => {
  if (!username.value || !password.value) {
    error.value = "Заполните все обязательные поля";
    return;
  }

  loading.value = true;
  error.value = null;
  success.value = false;

  try {
    await $fetch(baseURL + "/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
        role: role.value || undefined,
      }),
    });

    success.value = true;

    // Перенаправляем на страницу входа через 2 секунды
    setTimeout(() => {
      router.push("/login");
    }, 2000);
  } catch (e: any) {
    console.error("Ошибка регистрации:", e);
    error.value =
      e.data?.error || e.message || "Ошибка регистрации. Попробуйте еще раз.";
  } finally {
    loading.value = false;
  }
};
</script>
