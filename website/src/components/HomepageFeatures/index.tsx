import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Catches what costs you',
    description: (
      <>
        Missing memory limiters, hardcoded secrets, misordered processors —
        augur encodes hard-won operational knowledge into automated checks that
        run in milliseconds.
      </>
    ),
  },
  {
    title: 'Transparent Rego rules',
    description: (
      <>
        Every rule is a plain <code>.rego</code> file under <code>policy/</code>.
        Read them, override them, or write your own — no magic, no plugins, no
        config DSL to learn.
      </>
    ),
  },
  {
    title: 'CI-friendly from day one',
    description: (
      <>
        Text, JSON, and GitHub Actions annotation output. Non-zero exit on
        failure. Skip rules, merge custom policies, and promote warnings with{' '}
        <code>--strict</code>.
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
