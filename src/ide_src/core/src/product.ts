/** Human-readable name of the workbench, shown in the window title and About view. */
export const PRODUCT_NAME = 'FreeLunch IDE';

/** Reverse-DNS prefix for IDE preference keys, command ids, and contribution points. */
export const PRODUCT_NAMESPACE = 'freelunch';

/**
 * Label prefix for the Kubernetes labels FreeLunch stamps onto generated L2
 * resources. Roadmap section 6.3 requires `freelunch.io/product`,
 * `freelunch.io/workload`, and `freelunch.io/workload-id` to be stable, because
 * cost allocation joins on them.
 */
export const LABEL_PREFIX = 'freelunch.io';

/** Builds a fully-qualified FreeLunch resource label key. */
export function labelKey(name: string): string {
  if (name === '') {
    throw new Error('label name must not be empty');
  }
  return `${LABEL_PREFIX}/${name}`;
}
