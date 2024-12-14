export interface ReceiptItem {
  id: string;
  name: string;
  quantity: number;
  unitPrice: number;
  subtotal: number;
  totalPrice: number;
}

export interface Receipt {
  id: string;
  receiptNumber: string;
  date: string;
  customerName?: string;
  customerPhone?: string;
  paymentMethod: 'cash' | 'mpesa' | 'credit';
  mpesaReference?: string;
  items: ReceiptItem[];
  subtotal: number;
  tax: number;
  discount: number;
  totalAmount: number;
  businessName: string;
  businessPhone?: string;
  businessAddress?: string;
} 