import DashboardLayout from '../../dashboard/DashboardLayout';
import ReceiptMetrics from '../../dashboard/receipts/ReceiptMetrics';
import ReceiptList from '../../dashboard/receipts/ReceiptList';
import { useTranslation } from 'react-i18next';

export default function ReceiptManagement() {
  const { t } = useTranslation();

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <h1 className="text-2xl font-bold text-[#011627]">{t('receipts.title')}</h1>
        <ReceiptMetrics />
        <ReceiptList />
      </div>
    </DashboardLayout>
  );
} 