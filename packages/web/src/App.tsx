import { lazy, Suspense } from "react";
import { Routes, Route } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";

const Landing = lazy(() => import("@/pages/Landing"));
const Dashboard = lazy(() => import("@/pages/Dashboard"));
const Deposit = lazy(() => import("@/pages/Deposit"));
const Services = lazy(() => import("@/pages/Services"));
const RegionDetail = lazy(() => import("@/pages/RegionDetail"));
const Billing = lazy(() => import("@/pages/Billing"));
const BillDetail = lazy(() => import("@/pages/BillDetail"));
const Notifications = lazy(() => import("@/pages/Notifications"));

function PageFallback() {
  return (
    <div className="flex items-center justify-center min-h-[50vh]">
      <div className="text-text-secondary text-sm">Loading...</div>
    </div>
  );
}

export default function App() {
  return (
    <Suspense fallback={<PageFallback />}>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route element={<AppLayout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/deposit" element={<Deposit />} />
          <Route path="/services" element={<Services />} />
          <Route path="/services/:regionCode" element={<RegionDetail />} />
          <Route path="/billing" element={<Billing />} />
          <Route path="/billing/:billId" element={<BillDetail />} />
          <Route path="/notifications" element={<Notifications />} />
        </Route>
      </Routes>
    </Suspense>
  );
}
