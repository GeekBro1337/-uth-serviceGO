export const useApi = () => {
  const config = useRuntimeConfig();
  const baseURL = config.public.apiBase || "http://localhost:8999";

  const accessToken = useCookie("access_token", {
    maxAge: 60 * 15, // 15 минут
    sameSite: "lax",
    secure: false, // для localhost
    httpOnly: false, // нужно для чтения в JS
  });
  const refreshToken = useCookie("refresh_token", {
    maxAge: 60 * 60 * 24 * 7, // 7 дней
    sameSite: "lax",
    secure: false, // для localhost
    httpOnly: false, // нужно для чтения в JS
  });

  let isRefreshing = false;
  let refreshPromise: Promise<string> | null = null;

  // Авто-попытка подтянуть access из refresh cookie при инициализации (клиент)
  const ensureAccessToken = async () => {
    if (typeof window === "undefined") return; // только на клиенте
    if (!accessToken.value && !isRefreshing) {
      try {
        await refreshAccessToken();
      } catch {
        // игнорируем — пользователь может быть не авторизован
      }
    }
  };

  const refreshAccessToken = async (): Promise<string> => {
    if (refreshPromise) {
      return refreshPromise;
    }

    refreshPromise = (async () => {
      try {
        // Refresh токен автоматически отправляется в cookie (HttpOnly)
        const response: any = await $fetch(baseURL + "/refresh", {
          method: "POST",
          credentials: "include", // важно для отправки cookies
        });

        // В новой версии ответ содержит "token" вместо "access_token"
        const newToken: string | undefined = response.token || response.access_token;
        if (!newToken) {
          throw new Error("Access token not returned by refresh");
        }

        accessToken.value = newToken;
        return newToken;
      } catch (error) {
        // Если refresh не удался, очищаем токены
        accessToken.value = null;
        refreshToken.value = null;
        throw error;
      } finally {
        refreshPromise = null;
      }
    })();

    return refreshPromise;
  };

  const request = async (url: string, options: any = {}) => {
    const headers: any = { "Content-Type": "application/json" };

    // Если access отсутствует, пробуем обновить заранее
    if (!accessToken.value && !isRefreshing) {
      try {
        const newToken = await refreshAccessToken();
        if (newToken) {
          headers["Authorization"] = `Bearer ${newToken}`;
        }
      } catch (err) {
        // если не удалось, запрос пойдёт далее и вернёт 401 — обработаем ниже
      }
    }

    if (accessToken.value) {
      headers["Authorization"] = `Bearer ${accessToken.value}`;
    } else {
      console.warn("Access token отсутствует для запроса:", url);
    }

    try {
      const res = await $fetch(baseURL + url, { ...options, headers });
      return res;
    } catch (error: any) {
      console.error("Ошибка запроса:", url, error.status, error.message);
      // Если получили 401 и есть refresh token, пытаемся обновить
      if (error.status === 401 && !isRefreshing) {
        isRefreshing = true;
        try {
          const newAccessToken = await refreshAccessToken();
          // Повторяем запрос с новым токеном
          headers["Authorization"] = `Bearer ${newAccessToken}`;
          const res = await $fetch(baseURL + url, { ...options, headers });
          return res;
        } catch (refreshError) {
          // Если обновление не удалось, перенаправляем на логин
          // Используем window.location для надежного редиректа
          if (typeof window !== "undefined") {
            window.location.href = "/login";
          }
          throw refreshError;
        } finally {
          isRefreshing = false;
        }
      }
      throw error;
    }
  };

  const logout = async () => {
    try {
      // Вызываем серверный endpoint для удаления refresh токена из БД
      await $fetch(baseURL + "/logout", {
        method: "POST",
        credentials: "include", // важно для отправки cookies
      });
    } catch (error) {
      console.error("Ошибка при выходе:", error);
    } finally {
      // Очищаем токены на клиенте в любом случае
      accessToken.value = null;
      refreshToken.value = null;
    }
  };

  return {
    request,
    accessToken,
    refreshToken,
    logout,
    refreshAccessToken,
  };

  // Авто-попытка обновить access при инициализации (когда есть refresh в cookie)
  ensureAccessToken();
};
