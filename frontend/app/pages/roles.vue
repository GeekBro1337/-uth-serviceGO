<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">Role Management</h1>
    <input
      v-model="superadminCode"
      type="password"
      placeholder="SuperAdmin Code"
      class="input mb-4"
    />
    <button @click="loadRoles" class="btn-primary">Load Roles</button>

    <div
      v-for="r in roles"
      :key="r.role"
      class="border p-4 mt-4 rounded bg-gray-50"
    >
      <h2 class="font-bold text-lg">{{ r.role }}</h2>
      <p class="text-sm mb-2">{{ r.description }}</p>

      <div class="flex flex-wrap gap-2">
        <label
          v-for="perm in allPerms"
          :key="perm"
          class="flex items-center gap-2"
        >
          <input type="checkbox" v-model="r.permissions" :value="perm" />
          <span>{{ perm }}</span>
        </label>
      </div>
    </div>

    <button @click="save" class="btn-primary mt-6">Save Changes</button>
  </div>
</template>

<script setup lang="ts">
const { request } = useApi();
const superadminCode = ref("");
const roles = ref<any[]>([]);
const allPerms = ["read", "edit", "delete", "create"];

const loadRoles = async () => {
  roles.value = await request("/roles");
};

const save = async () => {
  await request("/roles/update", {
    method: "POST",
    body: JSON.stringify({
      superadmin_code: superadminCode.value,
      updates: roles.value.map((r) => ({
        role: r.role,
        permissions: r.permissions,
      })),
    }),
  });
  alert("Roles updated");
};
</script>
