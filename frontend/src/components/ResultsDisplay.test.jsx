import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';

import ResultsDisplay from './ResultsDisplay';

describe('ResultsDisplay', () => {

  it('should display the AI summary answer from props', () => {
    const mockResult = {
      ai_summary_answer: 'This is a test summary from the AI.'
    };

    render(<ResultsDisplay result={mockResult} />);

    const summaryElement = screen.getByText('This is a test summary from the AI.');
    
    expect(summaryElement).toBeInTheDocument();
  });

  it('should render the results header with icon and title', () => {
    const mockResult = {
      ai_summary_answer: 'Test summary'
    };

    render(<ResultsDisplay result={mockResult} />);

    const resultsIcon = screen.getByText('💡');
    const resultsTitle = screen.getByText('AI Summary');
    
    expect(resultsIcon).toBeInTheDocument();
    expect(resultsTitle).toBeInTheDocument();
  });

  it('should render the summary card with proper structure', () => {
    const mockResult = {
      ai_summary_answer: 'Test summary content'
    };

    render(<ResultsDisplay result={mockResult} />);

    const summaryCard = screen.getByText('Test summary content').closest('.summary-card');
    expect(summaryCard).toBeInTheDocument();
  });

  it('should render the AI badge in the summary footer', () => {
    const mockResult = {
      ai_summary_answer: 'Test summary'
    };

    render(<ResultsDisplay result={mockResult} />);

    const aiIcon = screen.getByText('🤖');
    const aiGeneratedText = screen.getByText('AI Generated');
    
    expect(aiIcon).toBeInTheDocument();
    expect(aiGeneratedText).toBeInTheDocument();
  });

  it('should render the complete component structure', () => {
    const mockResult = {
      ai_summary_answer: 'Complete test summary with multiple sentences to verify the full component structure is working correctly.'
    };

    render(<ResultsDisplay result={mockResult} />);

    // Check main container
    const resultsDisplay = screen.getByText('Complete test summary with multiple sentences to verify the full component structure is working correctly.').closest('.results-display');
    expect(resultsDisplay).toBeInTheDocument();

    // Check header elements
    expect(screen.getByText('💡')).toBeInTheDocument();
    expect(screen.getByText('AI Summary')).toBeInTheDocument();

    // Check content structure
    expect(screen.getByText('Complete test summary with multiple sentences to verify the full component structure is working correctly.')).toBeInTheDocument();

    // Check footer elements
    expect(screen.getByText('🤖')).toBeInTheDocument();
    expect(screen.getByText('AI Generated')).toBeInTheDocument();
  });

  it('should handle long summary text', () => {
    const longSummary = 'This is a very long summary that contains multiple sentences and should be displayed properly within the summary card. It should wrap correctly and maintain proper formatting. The text should be readable and well-structured.';
    
    const mockResult = {
      ai_summary_answer: longSummary
    };

    render(<ResultsDisplay result={mockResult} />);

    const summaryElement = screen.getByText(longSummary);
    expect(summaryElement).toBeInTheDocument();
  });

  it('should handle short summary text', () => {
    const shortSummary = 'Short answer.';
    
    const mockResult = {
      ai_summary_answer: shortSummary
    };

    render(<ResultsDisplay result={mockResult} />);

    const summaryElement = screen.getByText(shortSummary);
    expect(summaryElement).toBeInTheDocument();
  });

  it('should handle empty summary text', () => {
    const mockResult = {
      ai_summary_answer: ''
    };

    render(<ResultsDisplay result={mockResult} />);

    // The component should still render the structure even with empty content
    expect(screen.getByText('AI Summary')).toBeInTheDocument();
    expect(screen.getByText('🤖')).toBeInTheDocument();
    expect(screen.getByText('AI Generated')).toBeInTheDocument();
  });

  it('should have proper CSS classes for styling', () => {
    const mockResult = {
      ai_summary_answer: 'Test summary'
    };

    const { container } = render(<ResultsDisplay result={mockResult} />);

    // Check for main container class
    const resultsDisplay = container.querySelector('.results-display');
    expect(resultsDisplay).toBeInTheDocument();

    // Check for header classes
    const resultsHeader = container.querySelector('.results-header');
    expect(resultsHeader).toBeInTheDocument();

    // Check for content classes
    const resultsContent = container.querySelector('.results-content');
    expect(resultsContent).toBeInTheDocument();

    // Check for card classes
    const summaryCard = container.querySelector('.summary-card');
    expect(summaryCard).toBeInTheDocument();

    // Check for text classes
    const summaryText = container.querySelector('.summary-text');
    expect(summaryText).toBeInTheDocument();

    // Check for footer classes
    const summaryFooter = container.querySelector('.summary-footer');
    expect(summaryFooter).toBeInTheDocument();

    // Check for badge classes
    const aiBadge = container.querySelector('.ai-badge');
    expect(aiBadge).toBeInTheDocument();
  });
});