import { useState, useEffect } from 'react';
import { TrendingUp, DollarSign, CreditCard, ShoppingBag } from 'lucide-react';
import { SalesMetrics } from '../../utils/api';
import { formatCurrency } from '../../utils/formatters';
import { useTranslation } from 'react-i18next';
import { inventoryApi } from '../../utils/api';
export default function SalesInsights() {
  const { t } = useTranslation();
  const [metrics, setMetrics] = useState<SalesMetrics | null>(null);

  useEffect(() => {
    async function fetchMetrics() {
      try {
        const response = await inventoryApi.fetchSalesMetrics();
      
        if (response.success && response.data) {
          console.log("Metric data",response.data.topProducts[0]);
          setMetrics(response.data);
        } else {
          throw new Error(response.error || 'Failed to fetch sales metrics');
        }
      } catch (error) {
        console.error('Error fetching sales metrics:', error instanceof Error ? error.message : 'An unknown error occurred');
      }
    }

    fetchMetrics();
  }, []);

  if (!metrics) {
    return <div>Loading...</div>;
  }

  const stats = [
    {
      name: t('salesInsights.metrics.dailyRevenue'),
      value: formatCurrency(metrics.dailyRevenue),
      icon: DollarSign,
      change: '+4.75%',
      changeType: 'increase',
    },
    {
      name: t('salesInsights.metrics.weeklyRevenue'),
      value: formatCurrency(metrics.weeklyRevenue),
      icon: TrendingUp,
      change: '+54.02%',
      changeType: 'increase',
    },
    {
      name: t('salesInsights.metrics.monthlyRevenue'),
      value: formatCurrency(metrics.monthlyRevenue),
      icon: CreditCard,
      change: '+12.25%',
      changeType: 'increase',
    },
    {
      name: t('salesInsights.metrics.topProduct'),
      value: metrics.topProducts[0].product_name || 'N/iA',
      icon: ShoppingBag,
      change: metrics.topProducts[0]?.quantity 
        ? `${metrics.topProducts[0].quantity} ${t('salesInsights.units')}`
        : t('salesInsights.units') + ' 0',
      changeType: 'neutral',
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {stats.map((stat) => (
        <div key={stat.name} className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-500">{stat.name}</p>
              <p className="mt-1 text-xl font-semibold text-[#011627]">{stat.value}</p>
            </div>
            <div className="bg-[#2EC4B6]/10 rounded-full p-3">
              <stat.icon className="w-6 h-6 text-[#2EC4B6]" />
            </div>
          </div>
          <div className="mt-4">
            <div className={`inline-flex items-center text-sm ${
              stat.changeType === 'increase'
                ? 'text-green-600'
                : stat.changeType === 'decrease'
                ? 'text-red-600'
                : 'text-gray-500'
            }`}>
              <span>{stat.change}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}