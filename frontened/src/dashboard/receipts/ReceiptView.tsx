import { useState } from 'react';
import { motion } from 'framer-motion';
import { Download, X } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { formatDate, formatCurrency } from '../../utils/formatters';

interface ReceiptViewProps {
  receipt: {
    id: string;
    receiptNumber: string;
    date: string;
    customerName?: string;
    customerPhone?: string;
    paymentMethod: string;
    mpesaReference?: string;
    items: Array<{
      name: string;
      quantity: number;
      unitPrice: number;
      subtotal: number;
    }>;
    subtotal: number;
    tax: number;
    discount: number;
    total: number;
  };
  onClose: () => void;
}

export default function ReceiptView({ receipt, onClose }: ReceiptViewProps) {
  const { t } = useTranslation();
  const [isDownloading, setIsDownloading] = useState(false);

  const handleDownload = async () => {
    setIsDownloading(true);
    try {
      // Add PDF generation and download logic here
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate download
    } catch (error) {
      console.error('Failed to download receipt:', error);
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
    >
      <motion.div
        initial={{ scale: 0.95 }}
        animate={{ scale: 1 }}
        exit={{ scale: 0.95 }}
        className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
      >
        <div className="p-6 border-b border-gray-200">
          <div className="flex justify-between items-center">
            <h2 className="text-xl font-semibold text-[#011627]">
              {t('receipts.view.title')}
            </h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-500"
            >
              <X className="w-6 h-6" />
            </button>
          </div>
        </div>

        <div className="p-6 space-y-6">
          {/* Header Information */}
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-gray-500">{t('receipts.view.receiptNumber')}</span>
              <span className="font-medium">{receipt.receiptNumber}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">{t('receipts.view.date')}</span>
              <span>{formatDate(receipt.date)}</span>
            </div>
            {receipt.customerName && (
              <div className="flex justify-between">
                <span className="text-gray-500">{t('receipts.view.customer')}</span>
                <span>{receipt.customerName}</span>
              </div>
            )}
          </div>

          {/* Items Table */}
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">
                  {t('receipts.view.item')}
                </th>
                <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">
                  {t('receipts.view.quantity')}
                </th>
                <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">
                  {t('receipts.view.price')}
                </th>
                <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">
                  {t('receipts.view.subtotal')}
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {receipt.items.map((item, index) => (
                <tr key={index}>
                  <td className="px-4 py-2">{item.name}</td>
                  <td className="px-4 py-2 text-right">{item.quantity}</td>
                  <td className="px-4 py-2 text-right">{formatCurrency(item.unitPrice)}</td>
                  <td className="px-4 py-2 text-right">{formatCurrency(item.subtotal)}</td>
                </tr>
              ))}
            </tbody>
          </table>

          {/* Totals */}
          <div className="space-y-2 border-t pt-4">
            <div className="flex justify-between">
              <span className="text-gray-500">{t('receipts.view.subtotal')}</span>
              <span>{formatCurrency(receipt.subtotal)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">{t('receipts.view.tax')}</span>
              <span>{formatCurrency(receipt.tax)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">{t('receipts.view.discount')}</span>
              <span>{formatCurrency(receipt.discount)}</span>
            </div>
            <div className="flex justify-between font-semibold text-lg">
              <span>{t('receipts.view.total')}</span>
              <span>{formatCurrency(receipt.total)}</span>
            </div>
          </div>

          {/* Payment Information */}
          <div className="space-y-2 border-t pt-4">
            <div className="flex justify-between">
              <span className="text-gray-500">{t('receipts.view.paymentMethod')}</span>
              <span>{receipt.paymentMethod}</span>
            </div>
            {receipt.mpesaReference && (
              <div className="flex justify-between">
                <span className="text-gray-500">{t('receipts.view.mpesaReference')}</span>
                <span>{receipt.mpesaReference}</span>
              </div>
            )}
          </div>
        </div>

        {/* Footer Actions */}
        <div className="p-6 border-t border-gray-200">
          <button
            onClick={handleDownload}
            disabled={isDownloading}
            className="w-full flex items-center justify-center bg-[#2EC4B6] text-white py-2 rounded-lg hover:bg-[#28b0a3] transition-colors disabled:opacity-50"
          >
            <Download className="w-5 h-5 mr-2" />
            {isDownloading ? t('receipts.view.downloading') : t('receipts.view.download')}
          </button>
        </div>
      </motion.div>
    </motion.div>
  );
} 