import { describe, expect, it } from 'vitest';

import { LABEL_PREFIX, PRODUCT_NAME, labelKey } from './index';

describe('product constants', () => {
  it('exposes the workbench name', () => {
    expect(PRODUCT_NAME).toBe('FreeLunch IDE');
  });

  it('builds label keys under the freelunch.io prefix', () => {
    expect(labelKey('workload')).toBe(`${LABEL_PREFIX}/workload`);
  });

  it('rejects an empty label name', () => {
    expect(() => labelKey('')).toThrow(/must not be empty/);
  });
});
