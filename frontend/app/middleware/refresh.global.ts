export default defineNuxtRouteMiddleware(async () => {
  // Выполняем только на клиенте (refresh дергает /refresh с cookies)
  if (import.meta.server) return;

  const { accessToken, refreshAccessToken } = useApi();

  // Если access отсутствует, пробуем обновить его из refresh cookie
  if (!accessToken.value) {
    try {
      await refreshAccessToken();
    } catch {
      // игнорируем — пользователь может быть не авторизован
    }
  }
});

