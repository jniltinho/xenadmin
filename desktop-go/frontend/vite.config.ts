import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";

// --mode desktop (default) builds into internal/gui/dist, embedded in the binary; --mode web into ../web/dist.
export default defineConfig(({ mode }) => ({
  plugins: [vue(), tailwindcss()],
  base: "./",
  build: { target: "es2022", outDir: mode === "web" ? "../web/dist" : "../internal/gui/dist", emptyOutDir: true },
}));
