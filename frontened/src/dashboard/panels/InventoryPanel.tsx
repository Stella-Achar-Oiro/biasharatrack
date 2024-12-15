import { Package, AlertCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

export default function InventoryPanel() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-[#011627]">{t('inventoryPanel.title')}</h2>
        <Package className="h-6 w-6 text-[#2EC4B6]" />
      </div>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-500">{t('inventoryPanel.totalItems')}</span>
          <span className="text-lg font-medium">156</span>
        </div>
        <div className="flex items-center text-[#FF9F1C]">
          <AlertCircle className="h-5 w-5 mr-2" />
          <span className="text-sm">{t('inventoryPanel.lowStock', { count: 5 })}</span>
        </div>
        <button 
          onClick={() => navigate('/dashboard/inventory')}
          className="w-full bg-[#2EC4B6] text-white py-2 rounded-lg hover:bg-[#28b0a3] transition-colors"
        >
          {t('inventoryPanel.button')}
        </button>
      </div>
    </div>
  );
}