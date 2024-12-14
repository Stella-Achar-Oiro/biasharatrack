import { useState, useEffect} from 'react';
import { Search, Download, Eye, X } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { formatDate, formatCurrency } from '../../utils/formatters';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { Receipt } from '../../types/receipt';
import { inventoryApi } from '../../utils/api';
import jsPDF from 'jspdf';
import { useAuth } from '../../context/AuthContext';

export default function ReceiptList() {
  const { t } = useTranslation();
  const { user } = useAuth();
  const [searchTerm, setSearchTerm] = useState('');
  const [paymentFilter, setPaymentFilter] = useState('all');
  const [startDate, setStartDate] = useState<Date | null>(null);
  const [endDate, setEndDate] = useState<Date | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedReceipt, setSelectedReceipt] = useState<Receipt | null>(null);
  const receiptsPerPage = 5;

  useEffect(() => {
    const fetchReceipts = async () => {
      setLoading(true);
      const response = await inventoryApi.fetchReceipts();
      console.log("response====>",response);
      if (response.success && response.data) {
        setReceipts(response.data);
        setError(null);
      } else {
        setError(response.error ?? 'An error occurred');
        setReceipts([]);
      }
      setLoading(false);
    };

    fetchReceipts();
  }, []);

  const filteredReceipts = receipts.filter((receipt) => {
    const matchesSearch = 
      receipt.receiptNumber.toLowerCase().includes(searchTerm.toLowerCase()) ||
      receipt.customerName?.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesPayment = paymentFilter === 'all' || receipt.paymentMethod === paymentFilter;
    const matchesDate = (!startDate || new Date(receipt.date) >= startDate) &&
                       (!endDate || new Date(receipt.date) <= endDate);
    return matchesSearch && matchesPayment && matchesDate;
  });

  const indexOfLastReceipt = currentPage * receiptsPerPage;
  const indexOfFirstReceipt = indexOfLastReceipt - receiptsPerPage;
  const currentReceipts = filteredReceipts.slice(indexOfFirstReceipt, indexOfLastReceipt);
  const totalPages = Math.ceil(filteredReceipts.length / receiptsPerPage);

  const handleViewReceipt = (receipt: Receipt) => {
    setSelectedReceipt(receipt);
  };

  const handleCloseReceipt = () => {
    setSelectedReceipt(null);
  };

  const handleDownloadReceipt = (receipt: Receipt) => {
    // Create PDF document
    const doc = new jsPDF();
    
    // Get business details from auth context
    const businessName = user?.businessName || 'Business Name Not Set';
    const businessPhone = user?.telephone || 'Phone Not Set';
    const businessAddress = user?.location || 'Address Not Set';
    const customerName = receipt.customerName || 'Walk-in Customer';
    const tax = receipt.tax || 0;
    const discount = receipt.discount || 0;

    // Add business header
    doc.setFontSize(20);
    doc.setFont('helvetica', 'bold');
    doc.text(businessName, 105, 20, { align: 'center' });
    
    doc.setFontSize(12);
    doc.setFont('helvetica', 'normal');
    doc.text(`Tel: ${businessPhone}`, 105, 30, { align: 'center' });
    doc.text(businessAddress, 105, 40, { align: 'center' });

    // Add receipt details
    doc.setFont('helvetica', 'bold');
    doc.text('Receipt Details:', 20, 60);
    doc.setFont('helvetica', 'normal');
    doc.text(`Receipt #: ${receipt.receiptNumber}`, 20, 70);
    doc.text(`Date: ${formatDate(receipt.date)}`, 20, 80);
    doc.text(`Customer: ${customerName}`, 20, 90);
    doc.text(`Payment Method: ${receipt.paymentMethod.toUpperCase()}`, 20, 100);

    // Add items table
    doc.setFont('helvetica', 'bold');
    doc.text('Item', 20, 120);
    doc.text('Qty', 100, 120);
    doc.text('Price', 140, 120);
    doc.text('Total', 180, 120);
    
    let yPos = 130;
    doc.setFont('helvetica', 'normal');
    receipt.items.forEach(item => {
      doc.text(item.name, 20, yPos);
      doc.text(item.quantity.toString(), 100, yPos);
      doc.text(formatCurrency(item.unitPrice), 140, yPos);
      doc.text(formatCurrency(item.totalPrice), 180, yPos);
      yPos += 10;
    });

    // Add totals
    yPos += 10;
    doc.text(`Subtotal: ${formatCurrency(receipt.totalAmount)}`, 140, yPos);
    doc.text(`Tax: ${formatCurrency(tax)}`, 140, yPos + 10);
    doc.text(`Discount: ${formatCurrency(discount)}`, 140, yPos + 20);
    
    doc.setFont('helvetica', 'bold');
    doc.text(`Total: ${formatCurrency(receipt.totalAmount)}`, 140, yPos + 30);

    // Add footer
    doc.setFont('helvetica', 'normal');
    doc.setFontSize(10);
    doc.text('Thank you for your business!', 105, yPos + 50, { align: 'center' });

    // Save PDF
    doc.save(`receipt-${receipt.receiptNumber}.pdf`);
  };

  if (loading) {
    return <div className="p-6 text-center">Loading receipts...</div>;
  }

  if (error) {
    return <div className="p-6 text-center text-red-500">Error: {error}</div>;
  }

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-6 border-b border-gray-200">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder={t('receipts.search')}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#2EC4B6] focus:border-transparent"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
          <div className="flex flex-wrap gap-4">
            <select
              className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#2EC4B6] focus:border-transparent"
              value={paymentFilter}
              onChange={(e) => setPaymentFilter(e.target.value)}
            >
              <option value="all">{t('receipts.filters.all')}</option>
              <option value="cash">{t('receipts.filters.cash')}</option>
              <option value="mpesa">{t('receipts.filters.mpesa')}</option>
              <option value="credit">{t('receipts.filters.credit')}</option>
            </select>
            <div className="flex flex-col md:flex-row items-center gap-2">
              <DatePicker
                selected={startDate}
                onChange={setStartDate}
                placeholderText={t('receipts.startDate')}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#2EC4B6] focus:border-transparent"
              />
              <DatePicker
                selected={endDate}
                onChange={setEndDate}
                placeholderText={t('receipts.endDate')}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#2EC4B6] focus:border-transparent"
              />
            </div>
          </div>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {t('receipts.table.receiptNumber')}
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {t('receipts.table.date')}
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {t('receipts.table.customer')}
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {t('receipts.table.amount')}
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {t('receipts.table.payment')}
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                {t('receipts.table.actions')}
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {currentReceipts.map((receipt) => (
              <tr key={receipt.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {receipt.receiptNumber}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(receipt.date)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {receipt.customerName || '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {formatCurrency(receipt.totalAmount)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                    receipt.paymentMethod === 'cash'
                      ? 'bg-green-100 text-green-800'
                      : receipt.paymentMethod === 'mpesa'
                      ? 'bg-blue-100 text-blue-800'
                      : 'bg-yellow-100 text-yellow-800'
                  }`}>
                    {receipt.paymentMethod.toUpperCase()}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div className="flex justify-end space-x-2">
                    <button
                      onClick={() => handleViewReceipt(receipt)}
                      className="text-[#2EC4B6] hover:text-[#28b0a3]"
                      title={t('receipts.actions.view')}
                    >
                      <Eye className="w-5 h-5" />
                    </button>
                    <button
                      onClick={() => handleDownloadReceipt(receipt)}
                      className="text-[#FF9F1C] hover:text-[#f39200]"
                      title={t('receipts.actions.download')}
                    >
                      <Download className="w-5 h-5" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {filteredReceipts.length > receiptsPerPage && (
        <div className="flex justify-center items-center space-x-2 p-4 border-t">
          <button
            onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
            disabled={currentPage === 1}
            className="px-3 py-1 rounded-md border disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
          >
            {t('common.previous')}
          </button>
          <span className="text-sm text-gray-600">
            {t('common.page')} {currentPage} {t('common.of')} {totalPages}
          </span>
          <button
            onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
            disabled={currentPage === totalPages}
            className="px-3 py-1 rounded-md border disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
          >
            {t('common.next')}
          </button>
        </div>
      )}

      {/* Receipt Modal */}
      {selectedReceipt && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white w-[595px] h-[842px] p-8 rounded-lg shadow-lg relative overflow-y-auto">
            <button
              onClick={handleCloseReceipt}
              className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
            >
              <X className="w-6 h-6" />
            </button>
            
            <div className="text-center mb-6">
              <h2 className="text-xl font-bold">{selectedReceipt.businessName}</h2>
              {selectedReceipt.businessPhone && <p>Tel: {selectedReceipt.businessPhone}</p>}
              {selectedReceipt.businessAddress && <p>{selectedReceipt.businessAddress}</p>}
            </div>

            <div className="mb-6">
              <p><strong>Receipt #:</strong> {selectedReceipt.receiptNumber}</p>
              <p><strong>Date:</strong> {formatDate(selectedReceipt.date)}</p>
              <p><strong>Customer:</strong> {selectedReceipt.customerName || 'Walk-in Customer'}</p>
              <p><strong>Payment Method:</strong> {selectedReceipt.paymentMethod.toUpperCase()}</p>
            </div>

            <table className="w-full mb-6">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2">Item</th>
                  <th className="text-right py-2">Qty</th>
                  <th className="text-right py-2">Price</th>
                  <th className="text-right py-2">Total</th>
                </tr>
              </thead>
              <tbody>
                {selectedReceipt.items.map((item, index) => (
                  <tr key={index} className="border-b">
                    <td className="py-2">{item.name}</td>
                    <td className="text-right py-2">{item.quantity}</td>
                    <td className="text-right py-2">{formatCurrency(item.unitPrice)}</td>
                    <td className="text-right py-2">{formatCurrency(item.totalPrice)}</td>
                  </tr>
                ))}
              </tbody>
            </table>

            <div className="text-right">
              <p><strong>Subtotal:</strong> {formatCurrency(selectedReceipt.totalAmount)}</p>
              <p><strong>Tax:</strong> {formatCurrency(selectedReceipt.tax)}</p>
              <p><strong>Discount:</strong> {formatCurrency(selectedReceipt.discount)}</p>
              <p className="text-xl font-bold mt-2">
                <strong>Total:</strong> {formatCurrency(selectedReceipt.totalAmount)}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Debug logs
      <div className="mt-4 text-sm text-gray-500">
        <p>Filtered Receipts: {filteredReceipts.length}</p>
        <p>Current Page: {currentPage}</p>
        <p>Current Receipts: {currentReceipts.length}</p>
      </div> */}
    </div>
  );
}