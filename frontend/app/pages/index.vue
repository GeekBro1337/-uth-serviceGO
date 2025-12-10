<template>
  <div
    class="min-h-screen bg-gradient-to-br from-gray-50 via-white to-blue-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900"
  >
    <div class="container mx-auto px-4 py-12">
      <!-- Hero секция для неавторизованных пользователей -->
      <div v-if="!user && !loading" class="max-w-4xl mx-auto text-center mb-12">
        <h1 class="text-5xl font-bold text-gray-900 dark:text-white mb-4">
          Добро пожаловать в Auth Service
        </h1>
        <p class="text-xl text-gray-600 dark:text-gray-300 mb-8">
          Система авторизации с поддержкой ролей и прав доступа
        </p>
        <div class="flex gap-4 justify-center">
          <UButton
            to="/login"
            size="xl"
            icon="i-heroicons-arrow-right-on-rectangle"
          >
            Войти
          </UButton>
          <UButton
            to="/register"
            size="xl"
            variant="outline"
            icon="i-heroicons-user-plus"
          >
            Зарегистрироваться
          </UButton>
        </div>
      </div>

      <!-- Контент для авторизованных пользователей -->
      <div v-if="user" class="max-w-4xl mx-auto">
        <UCard class="shadow-xl">
          <template #header>
            <div class="flex items-center gap-3">
              <UIcon
                name="i-heroicons-user-circle"
                class="w-8 h-8 text-primary-500"
              />
              <div>
                <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
                  Добро пожаловать, {{ user.username }}!
                </h1>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Ваш профиль
                </p>
              </div>
            </div>
          </template>

          <div class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <UCard>
                <div class="space-y-2">
                  <div
                    class="flex items-center gap-2 text-gray-600 dark:text-gray-400"
                  >
                    <UIcon name="i-heroicons-user" class="w-5 h-5" />
                    <span class="text-sm font-medium">Имя пользователя</span>
                  </div>
                  <p
                    class="text-lg font-semibold text-gray-900 dark:text-white"
                  >
                    {{ user.username }}
                  </p>
                </div>
              </UCard>

              <UCard>
                <div class="space-y-2">
                  <div
                    class="flex items-center gap-2 text-gray-600 dark:text-gray-400"
                  >
                    <UIcon name="i-heroicons-shield-check" class="w-5 h-5" />
                    <span class="text-sm font-medium">Роль</span>
                  </div>
                  <UBadge
                    :color="getRoleColor(user.role)"
                    variant="soft"
                    size="lg"
                  >
                    {{ user.role }}
                  </UBadge>
                </div>
              </UCard>
            </div>

            <div class="flex gap-4 pt-4">
              <UButton
                to="/me"
                icon="i-heroicons-user-circle"
                variant="outline"
              >
                Подробнее о профиле
              </UButton>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Состояние загрузки -->
      <div v-if="loading" class="w-full mx-auto">
        <UCard>
          <div class="text-center py-12">
            <UIcon
              name="i-heroicons-arrow-path"
              class="w-12 h-12 animate-spin text-primary-500 mx-auto mb-4"
            />
            <p class="text-gray-600 dark:text-gray-400">Загрузка...</p>
          </div>
        </UCard>
      </div>

      <!-- Ошибка -->
      <div v-if="error && !loading" class="max-w-2xl mx-auto mt-4">
        <UAlert
          color="error"
          variant="soft"
          :title="error"
          icon="i-heroicons-exclamation-triangle"
        />
      </div>

      <!-- Информационные карточки для всех пользователей -->
      <div class="max-w-6xl mx-auto mt-16">
        <h2
          class="text-3xl font-bold text-center text-gray-900 dark:text-white mb-8"
        >
          Возможности системы
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <UCard class="hover:shadow-lg transition-shadow">
            <div class="text-center">
              <UIcon
                name="i-heroicons-shield-check"
                class="w-12 h-12 text-primary-500 mx-auto mb-4"
              />
              <h3
                class="text-xl font-semibold text-gray-900 dark:text-white mb-2"
              >
                Безопасность
              </h3>
              <p class="text-gray-600 dark:text-gray-400">
                JWT токены с автоматическим обновлением для безопасной
                авторизации
              </p>
            </div>
          </UCard>

          <UCard class="hover:shadow-lg transition-shadow">
            <div class="text-center">
              <UIcon
                name="i-heroicons-user-group"
                class="w-12 h-12 text-primary-500 mx-auto mb-4"
              />
              <h3
                class="text-xl font-semibold text-gray-900 dark:text-white mb-2"
              >
                Роли и права
              </h3>
              <p class="text-gray-600 dark:text-gray-400">
                Гибкая система ролей с настраиваемыми правами доступа
              </p>
            </div>
          </UCard>

          <UCard class="hover:shadow-lg transition-shadow">
            <div class="text-center">
              <UIcon
                name="i-heroicons-bolt"
                class="w-12 h-12 text-primary-500 mx-auto mb-4"
              />
              <h3
                class="text-xl font-semibold text-gray-900 dark:text-white mb-2"
              >
                Производительность
              </h3>
              <p class="text-gray-600 dark:text-gray-400">
                Быстрая и надежная работа системы авторизации
              </p>
            </div>
          </UCard>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { request, accessToken } = useApi();

const user = ref<{ username: string; role: string } | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

const getRoleColor = (role: string) => {
  const colors: Record<string, "primary" | "success" | "warning" | "error"> = {
    admin: "error",
    user: "primary",
    auditor: "warning",
  };
  return colors[role] || "primary";
};

const loadUserData = async () => {
  // Если нет токена, просто не загружаем данные пользователя
  if (!accessToken.value) {
    loading.value = false;
    return;
  }

  try {
    loading.value = true;
    error.value = null;

    const data: any = await request("/protected/me");
    user.value = data;
  } catch (e: any) {
    // Если 401 - просто не показываем пользователя, это нормально для неавторизованных
    if (e.status !== 401) {
      console.error("Ошибка загрузки данных:", e);
      error.value =
        e.data?.error ||
        e.message ||
        "Не удалось загрузить данные пользователя";
    }
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadUserData();
});

// Следим за изменениями токена
watch(accessToken, (newToken) => {
  if (newToken) {
    loadUserData();
  } else {
    user.value = null;
  }
});
</script>
