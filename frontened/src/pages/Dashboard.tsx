import DashboardLayout from '../dashboard/DashboardLayout';
import InventoryPanel from '../dashboard/panels/InventoryPanel';
import SalesPanel from '../dashboard/panels/SalesPanel';
import FinancialPanel from '../dashboard/panels/FinancialPanel';
import CreditPanel from '../dashboard/panels/CreditPanel';
import AnalyticsPanel from '../dashboard/panels/AnalyticsPanel';
import TutorialPanel from '../dashboard/panels/TutorialPanel';

export default function Dashboard() {
  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <InventoryPanel />
          <SalesPanel />
          <FinancialPanel />
          <CreditPanel />
          <AnalyticsPanel />
          <TutorialPanel />
        </div>
      </div>
    </DashboardLayout>
  );
}