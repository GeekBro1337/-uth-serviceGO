<template>
  <header class="sticky top-0 z-50 w-full border-b bg-white/80 dark:bg-gray-900/80 backdrop-blur-md">
    <div class="container mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <!-- Левая часть - Информация об авторизации -->
        <div class="flex items-center gap-4">
          <NuxtLink to="/" class="flex items-center gap-2 hover:opacity-80 transition-opacity">
            <UIcon name="i-heroicons-shield-check" class="w-6 h-6 text-primary-500" />
            <span class="font-bold text-lg text-gray-900 dark:text-white">Auth Service</span>
          </NuxtLink>

          <div v-if="user" class="flex items-center gap-2 ml-4">
            <UIcon name="i-heroicons-user-circle" class="w-5 h-5 text-gray-600 dark:text-gray-400" />
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ user.username }}
            </span>
            <UBadge :color="getRoleColor(user.role)" variant="soft" size="xs">
              {{ user.role }}
            </UBadge>
          </div>

          <div v-else class="flex items-center gap-2 ml-4">
            <UIcon name="i-heroicons-user-circle" class="w-5 h-5 text-gray-400" />
            <span class="text-sm text-gray-500 dark:text-gray-400">Не авторизован</span>
          </div>
        </div>

        <!-- Правая часть - Кнопки действий -->
        <div class="flex items-center gap-2">
          <UColorModeButton />

          <template v-if="user">
            <UButton
              to="/me"
              variant="ghost"
              icon="i-heroicons-user"
              aria-label="Профиль"
            >
              Профиль
            </UButton>
            <UButton
              @click="handleLogout"
              color="error"
              variant="ghost"
              icon="i-heroicons-arrow-right-on-rectangle"
            >
              Выйти
            </UButton>
          </template>

          <template v-else>
            <UButton
              to="/login"
              variant="ghost"
              icon="i-heroicons-arrow-right-on-rectangle"
            >
              Войти
            </UButton>
            <UButton
              to="/register"
              icon="i-heroicons-user-plus"
            >
              Регистрация
            </UButton>
          </template>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const { accessToken, logout: apiLogout } = useApi();
const router = useRouter();

const user = ref<{ username: string; role: string } | null>(null);
const loading = ref(true);

const getRoleColor = (role: string) => {
  const colors: Record<string, "primary" | "success" | "warning" | "error"> = {
    admin: "error",
    user: "primary",
    auditor: "warning",
  };
  return colors[role] || "primary";
};

const loadUser = async () => {
  if (!accessToken.value) {
    loading.value = false;
    return;
  }

  try {
    const { request } = useApi();
    const data: any = await request("/protected/me");
    user.value = data;
  } catch (e: any) {
    // Если ошибка авторизации, просто не показываем пользователя
    if (e.status !== 401) {
      console.error("Ошибка загрузки пользователя:", e);
    }
  } finally {
    loading.value = false;
  }
};

const handleLogout = async () => {
  await apiLogout();
  user.value = null;
  router.push("/");
};

// Загружаем данные пользователя при монтировании
onMounted(() => {
  loadUser();
});

// Слушаем изменения токена
watch(accessToken, (newToken) => {
  if (newToken) {
    loadUser();
  } else {
    user.value = null;
  }
});
</script>

