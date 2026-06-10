/// <reference types="vitest/globals" />
// 让 tsc 识别 @testing-library/jest-dom 扩展的断言匹配器（toBeInTheDocument 等）。
// 运行时由 vitest.setup.ts 的 `import "@testing-library/jest-dom/vitest"` 注入。
import "@testing-library/jest-dom/vitest";
