import DashboardLayout from '../../dashboard/DashboardLayout';
import AddProduct from '../../dashboard/inventory/AddProduct';
import ProductList from '../../dashboard/inventory/ProductList';
import LowStockAlerts from '../../dashboard/inventory/LowStockAlerts';
import { useTranslation } from 'react-i18next';

export default function InventoryManagement() {
  const { t } = useTranslation();

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <h1 className="text-2xl font-bold text-[#011627]">{t('features.items.inventory.title')}</h1>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <AddProduct />
          <LowStockAlerts />
        </div>
        <ProductList />
      </div>
    </DashboardLayout>
  );
}