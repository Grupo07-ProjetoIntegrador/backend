import { createBrowserRouter } from "react-router";
import Root from "./Root";
import { TrainingManagement } from "./components/TrainingManagement";
import { DashboardAnalitico } from "./components/DashboardAnalitico";

export const router = createBrowserRouter([
  {
    path: "/",
    Component: Root,
    children: [
      { index: true, Component: DashboardAnalitico },
      { path: "treinamentos", Component: TrainingManagement },
    ],
  },
]);