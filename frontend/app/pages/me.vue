<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-white to-blue-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
    <div class="container mx-auto px-4 py-12">
      <div class="max-w-4xl mx-auto">
        <!-- Заголовок -->
        <div class="mb-8">
          <h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-2">
            Личный кабинет
          </h1>
          <p class="text-gray-600 dark:text-gray-400">
            Управление профилем и настройками
          </p>
        </div>

        <!-- Состояние загрузки -->
        <div v-if="loading" class="text-center py-12">
          <UIcon name="i-heroicons-arrow-path" class="w-12 h-12 animate-spin text-primary-500 mx-auto mb-4" />
          <p class="text-gray-600 dark:text-gray-400">Загрузка данных...</p>
        </div>

        <!-- Ошибка -->
        <UAlert
          v-else-if="error"
          color="error"
          variant="soft"
          :title="error"
          icon="i-heroicons-exclamation-triangle"
          class="mb-6"
        />

        <!-- Контент -->
        <div v-else-if="user" class="space-y-6">
          <!-- Основная информация -->
          <UCard class="shadow-xl">
            <template #header>
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                  <UIcon name="i-heroicons-user-circle" class="w-8 h-8 text-primary-600 dark:text-primary-400" />
                </div>
                <div>
                  <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                    {{ user.username }}
                  </h2>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    Информация о профиле
                  </p>
                </div>
              </div>
            </template>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Username -->
              <div class="space-y-2">
                <div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
                  <UIcon name="i-heroicons-user" class="w-5 h-5" />
                  <span class="text-sm font-medium">Имя пользователя</span>
                </div>
                <p class="text-lg font-semibold text-gray-900 dark:text-white">
                  {{ user.username }}
                </p>
              </div>

              <!-- Role -->
              <div class="space-y-2">
                <div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
                  <UIcon name="i-heroicons-shield-check" class="w-5 h-5" />
                  <span class="text-sm font-medium">Роль</span>
                </div>
                <UBadge :color="getRoleColor(user.role)" variant="soft" size="lg">
                  {{ user.role }}
                </UBadge>
              </div>
            </div>
          </UCard>

          <!-- Права доступа -->
          <UCard class="shadow-xl">
            <template #header>
              <div class="flex items-center gap-3">
                <UIcon name="i-heroicons-key" class="w-6 h-6 text-primary-500" />
                <div>
                  <h3 class="text-xl font-bold text-gray-900 dark:text-white">
                    Права доступа
                  </h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    Разрешения, доступные вашей роли
                  </p>
                </div>
              </div>
            </template>

            <template v-if="user.permissions && Array.isArray(user.permissions) && user.permissions.length > 0">
              <div class="space-y-3">
                <div
                  v-for="permission in user.permissions"
                  :key="permission"
                  class="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors"
                >
                  <UIcon
                    :name="getPermissionIcon(permission)"
                    class="w-5 h-5"
                    :class="getPermissionColor(permission)"
                  />
                  <span class="font-medium text-gray-900 dark:text-white">
                    {{ getPermissionLabel(permission) }}
                  </span>
                  <UBadge
                    :color="getPermissionBadgeColor(permission)"
                    variant="soft"
                    size="xs"
                    class="ml-auto"
                  >
                    {{ permission }}
                  </UBadge>
                </div>
              </div>
            </template>

            <template v-else>
              <div class="text-center py-8">
                <UIcon name="i-heroicons-lock-closed" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
                <p class="text-gray-600 dark:text-gray-400 mb-2">
                  У вашей роли нет назначенных прав доступа
                </p>
                <p class="text-sm text-gray-500 dark:text-gray-500">
                  Роль: {{ user.role }} | Права: {{ user.permissions ? JSON.stringify(user.permissions) : 'не загружены' }}
                </p>
              </div>
            </template>
          </UCard>

          <!-- Дополнительная информация -->
          <UCard class="shadow-xl">
            <template #header>
              <div class="flex items-center gap-3">
                <UIcon name="i-heroicons-information-circle" class="w-6 h-6 text-primary-500" />
                <h3 class="text-xl font-bold text-gray-900 dark:text-white">
                  Дополнительная информация
                </h3>
              </div>
            </template>

            <div class="space-y-4">
              <div class="flex items-start gap-3">
                <UIcon name="i-heroicons-shield-exclamation" class="w-5 h-5 text-gray-400 mt-0.5" />
                <div>
                  <p class="font-medium text-gray-900 dark:text-white mb-1">
                    Безопасность
                  </p>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    Ваши данные защищены с помощью JWT токенов. Access токен обновляется автоматически каждые 15 минут.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3">
                <UIcon name="i-heroicons-cog-6-tooth" class="w-5 h-5 text-gray-400 mt-0.5" />
                <div>
                  <p class="font-medium text-gray-900 dark:text-white mb-1">
                    Управление правами
                  </p>
                  <p class="text-sm text-gray-600 dark:text-gray-400">
                    Права доступа назначаются администратором системы в зависимости от вашей роли.
                  </p>
                </div>
              </div>
            </div>
          </UCard>
        </div>

        <!-- Не авторизован -->
        <UCard v-else class="shadow-xl">
          <div class="text-center py-12">
            <UIcon name="i-heroicons-exclamation-triangle" class="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">
              Доступ ограничен
            </h3>
            <p class="text-gray-600 dark:text-gray-400 mb-6">
              Для просмотра личного кабинета необходимо войти в систему
            </p>
            <UButton to="/login" icon="i-heroicons-arrow-right-on-rectangle">
              Войти
            </UButton>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { request } = useApi();
const router = useRouter();

const user = ref<{ username: string; role: string; permissions: string[] } | null>(null);
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

const getPermissionIcon = (permission: string) => {
  const icons: Record<string, string> = {
    read: "i-heroicons-eye",
    edit: "i-heroicons-pencil",
    delete: "i-heroicons-trash",
    create: "i-heroicons-plus-circle",
  };
  return icons[permission] || "i-heroicons-key";
};

const getPermissionLabel = (permission: string) => {
  const labels: Record<string, string> = {
    read: "Чтение",
    edit: "Редактирование",
    delete: "Удаление",
    create: "Создание",
  };
  return labels[permission] || permission;
};

const getPermissionColor = (permission: string) => {
  const colors: Record<string, string> = {
    read: "text-blue-500",
    edit: "text-yellow-500",
    delete: "text-red-500",
    create: "text-green-500",
  };
  return colors[permission] || "text-gray-500";
};

const getPermissionBadgeColor = (permission: string) => {
  const colors: Record<string, "primary" | "success" | "warning" | "error"> = {
    read: "primary",
    edit: "warning",
    delete: "error",
    create: "success",
  };
  return colors[permission] || "primary";
};

const loadUserData = async () => {
  try {
    loading.value = true;
    error.value = null;
    const data: any = await request("/protected/me");
    
    // Убеждаемся, что permissions всегда массив
    if (data && !Array.isArray(data.permissions)) {
      data.permissions = data.permissions || [];
    }
    
    console.log("Данные пользователя:", data);
    console.log("Права доступа:", data.permissions);
    
    user.value = data;
  } catch (e: any) {
    console.error("Ошибка загрузки данных:", e);
    error.value = e.data?.error || e.message || "Не удалось загрузить данные пользователя";
    if (e.status === 401) {
      router.push("/login");
    }
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadUserData();
});
</script>
