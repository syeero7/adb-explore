type Routes = "setup" | "explore";

export const router = $state<{ current: Routes }>({ current: "setup" });
