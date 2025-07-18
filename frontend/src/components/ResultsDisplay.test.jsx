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

});