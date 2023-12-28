import type { AppInterfaces } from '@shared/interfaces/appInterfaces';

import HostComponent from '../host/host';

type HomeProps = {
  hosts: AppInterfaces.Host[];
};

function HomeComponent(props: HomeProps) {
  const { hosts } = props;

  return (
    <div className=''>
      <div className='grid grid-cols-1 tablet:grid-cols-2 laptop:grid-cols-3 desktop:grid-cols-4 tv:grid-cols-5'>
        {hosts?.map((host, index) => (
          <HostComponent key={`${host.name}${index}`} host={host} />
        ))}
      </div>
    </div>
  );
}

export default HomeComponent;
