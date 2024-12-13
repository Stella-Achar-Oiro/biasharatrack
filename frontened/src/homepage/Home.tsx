import Hero from './components/Hero';
import Features from './components/Features';
import Impact from './components/Impact';
import CallToAction from './components/CallToAction';
import Footer from './components/Footer';
import React from 'react';

export default function Home() {
  return (
    <div className="min-h-screen">
      <Hero />
      <Features />
      <Impact />
      <CallToAction />
      <Footer />
    </div>
  );
}