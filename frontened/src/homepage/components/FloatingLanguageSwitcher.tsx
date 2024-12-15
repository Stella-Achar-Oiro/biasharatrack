import React, { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';

export default function FloatingLanguageSwitcher() {
  const { i18n } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const languages = [
    { code: 'en', label: 'English', shortLabel: 'EN' },
    { code: 'sw', label: 'Kiswahili', shortLabel: 'SW' }
  ];

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const changeLanguage = (langCode: string) => {
    i18n.changeLanguage(langCode);
    setIsOpen(false);
  };

  const getCurrentLanguage = () => {
    const currentLang = languages.find(lang => lang.code === i18n.language);
    return currentLang || languages[0];
  };

  return (
    <div 
      className="fixed top-4 right-12 xs:right-14 sm:right-16 md:right-20 lg:right-24 z-50" 
      ref={dropdownRef}
    >
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-1 xs:gap-2 px-2 xs:px-3 sm:px-4 py-1 xs:py-1.5 rounded-full bg-[#2EC4B6] text-white hover:bg-[#25a093] transition-colors shadow-lg text-xs xs:text-sm sm:text-base"
        aria-label="Toggle language"
        aria-expanded={isOpen}
        aria-haspopup="true"
      >
        <Globe className="w-3 h-3 xs:w-4 xs:h-4 sm:w-5 sm:h-5" />
        {/* Hide full text on mobile, show short form */}
        <span className="hidden sm:inline">
          {getCurrentLanguage().label}
        </span>
        {/* Show short form on mobile only */}
        <span className="inline sm:hidden">
          {getCurrentLanguage().shortLabel}
        </span>
        <span 
          className={`transition-transform duration-200 ${isOpen ? 'rotate-180' : ''} text-[10px] xs:text-xs sm:text-sm`}
          aria-hidden="true"
        >
          ▼
        </span>
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-1 w-24 xs:w-28 sm:w-36 bg-white rounded-lg shadow-lg overflow-hidden">
          {languages.map((lang) => (
            <button
              key={lang.code}
              onClick={() => changeLanguage(lang.code)}
              className={`w-full px-2 xs:px-3 sm:px-4 py-1.5 xs:py-1 text-left hover:bg-gray-100 transition-colors text-xs xs:text-sm sm:text-base ${
                i18n.language === lang.code 
                  ? 'bg-gray-50 text-[#2EC4B6]' 
                  : 'text-gray-700'
              }`}
            >
              <div className="flex items-center gap-1 xs:gap-2">
                {/* Show both short and full labels with responsive visibility */}
                <span className="hidden sm:inline">{lang.label}</span>
                <span className="inline sm:hidden">{lang.shortLabel}</span>
                {i18n.language === lang.code && (
                  <span className="text-[#2EC4B6] ml-auto">✓</span>
                )}
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}