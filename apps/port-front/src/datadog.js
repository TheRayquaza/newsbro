import { datadogRum } from '@datadog/browser-rum';
import { reactPlugin } from '@datadog/browser-rum-react';

export function initDatadog() {
  datadogRum.init({
      applicationId: '26ebd80b-edb8-497c-8dcf-7c083ebdfa99',
      clientToken: 'pub293c79f261a3e60c25b277a0314945f7',
      site: 'us5.datadoghq.com',
      service:'port-front',
      env: 'prod',
    
      // Specify a version number to identify the deployed version of your application in Datadog
      version: __APP_VERSION__,
      sessionSampleRate:  100,
      sessionReplaySampleRate: 100,
      defaultPrivacyLevel: 'allow',
      plugins: [reactPlugin({ router: true })],
  });
}
