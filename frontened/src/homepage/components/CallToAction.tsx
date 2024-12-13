import React from 'react';
import { ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { colors } from '../../utils/colors';
import { useTranslation } from 'react-i18next';

export default function CallToAction() {
  const { t } = useTranslation();

  return (
    <div className="bg-[#FDFFFC] py-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 className="text-3xl md:text-4xl font-bold text-[#011627] mb-6">
          {t('cta.title')}
        </h2>
        <p className="text-[#011627]/70 text-lg mb-8 max-w-2xl mx-auto">
          {t('cta.subtitle')}
        </p>
        <Link
          to="/signup"
          className="group inline-flex items-center gap-2 bg-[#E71D36] hover:bg-[#c91126] text-[#FDFFFC] px-8 py-4 rounded-full font-semibold text-lg transition-all"
        >
          {t('cta.button')}
          <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
        </Link>
      </div>
    </div>
  );
}