import { DollarSign, FileText, TrendingUp, Clock } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { formatCurrency } from '../../utils/formatters';

export default function ReceiptMetrics() {
  const { t } = useTranslation();

  // Replace with actual API data
  const metrics = {
    totalRevenue: 125000,
    totalReceipts: 45,
    averageAmount: 2777.78,
    recentReceipts: 12
  };

  const cards = [
    {
      title: t('receipts.metrics.totalRevenue'),
      value: formatCurrency(metrics.totalRevenue),
      icon: DollarSign,
      color: 'text-[#2EC4B6]',
      bgColor: 'bg-[#2EC4B6]/10'
    },
    {
      title: t('receipts.metrics.totalReceipts'),
      value: metrics.totalReceipts,
      icon: FileText,
      color: 'text-[#FF9F1C]',
      bgColor: 'bg-[#FF9F1C]/10'
    },
    {
      title: t('receipts.metrics.averageAmount'),
      value: formatCurrency(metrics.averageAmount),
      icon: TrendingUp,
      color: 'text-[#E71D36]',
      bgColor: 'bg-[#E71D36]/10'
    },
    {
      title: t('receipts.metrics.recentReceipts'),
      value: metrics.recentReceipts,
      icon: Clock,
      color: 'text-[#2EC4B6]',
      bgColor: 'bg-[#2EC4B6]/10'
    }
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {cards.map((card, index) => (
        <div key={index} className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-500">{card.title}</p>
              <p className="mt-1 text-xl font-semibold text-[#011627]">{card.value}</p>
            </div>
            <div className={`${card.bgColor} rounded-full p-3`}>
              <card.icon className={`w-6 h-6 ${card.color}`} />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
} 