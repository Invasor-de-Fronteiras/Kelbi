import { BoostModeTypes, LastAuthResult, SignResult } from './launcher';

if (process.env.NODE_ENV === 'development') {
  let lastAuthResult: LastAuthResult = LastAuthResult.None;
  let signResult: SignResult = SignResult.None;

  const username = 'test';
  const password = 'password';

  window.external.playSound = (): void => {
    console.log('playSound');
  };

  window.external.beginDrag = (_active: boolean): void => {
    console.log('beginDrag');
  };

  window.external.openBrowser = (_url: string): void => {
    console.log('openBrowser');
  };

  window.external.restartMhf = (): void => {
    console.log('restartMhf');
  };

  window.external.minimizeWindow = (): void => {
    console.log('minimizeWindow');
  };

  window.external.openMhlConfig = (): void => {
    console.log('openMhlConfig');
  };

  window.external.closeWindow = (): void => {
    console.log('closeWindow');
  };

  window.external.getUpdatePercentageTotal = () => {
    return 100;
  };

  window.external.getAccountRights = (): string => {
    // console.log('getAccountRights');
    return;
  };
  window.external.getMhfBootMode = (): BoostModeTypes => {
    // console.log('getMhfBootMode');
    return '_MHF_NORMAL';
  };

  window.external.getIniLastServerIndex = (): unknown => {
    // console.log('getIniLastServerIndex');
    return;
  };

  window.external.getServerListXml = (): string => {
    // console.log('getServerListXml');
    return `<?xml version="1.0"?><server_groups><group idx='0' nam='Erupe' ip='0.0.0.0' port="53312"/></server_groups>`;
  };

  window.external.getMhfMutexNumber = (): number => {
    // console.log('getMhfMutexNumber');
    return 0;
  };

  window.external.getUserId = (): string => {
    // console.log('getUserId');
    return username;
  };

  window.external.getPassword = (): string => {
    // console.log('getPassword');
    return password;
  };

  window.external.getLastAuthResult = (): LastAuthResult => {
    return lastAuthResult;
  };

  window.external.getSignResult = (): SignResult => {
    return signResult;
  };

  window.external.isEnableSessionId = (): unknown => {
    return 0;
  };

  window.external.startUpdate = () => {
    return true
  }

  window.external.getCharacterInfo = (): string => {
    return `<?xml version='1.0' encoding='shift_jis'?><CharacterInfo defaultUid=''><Character name='Slayer' uid='411111' weapon='???' HR='7' GR='13' lastLogin='1653524546' sex='F' />
        <Character name='Jesus' uid='911111' weapon='???' HR='99' GR='99' lastLogin='1653527484' sex='F' />
        <Character name='Fenda' uid='A11111' weapon='?????' HR='7' GR='0' lastLogin='1653566230' sex='F' />
    </CharacterInfo>
    `;
  };

  window.external.extractLog = (): unknown => {
    return null;
  };

  window.external.getUpdateStatus = (): number => {
    return 0;
  };

  window.external.getLauncherReturnCode = (): 'NORMAL' => {
    return 'NORMAL';
  };

  window.external.selectCharacter = (_charUid: string, _charUid1: string): void => {
    console.log('selectCharacter');
  };

  window.external.exitLauncher = (): void => {
    console.log('exitLauncher');
  };

  window.external.loginCog = (
    usernameInput: string,
    passwordInput: string,
    _confirmPassword: string,
  ): void => {
    lastAuthResult = LastAuthResult.InLoading;
    setTimeout(() => {
      const isCreatingChar = usernameInput.endsWith('+');
      if (isCreatingChar || (passwordInput === password && usernameInput === username)) {
        lastAuthResult = LastAuthResult.AuthSuccess;
        signResult = SignResult.SignSuccess;
      } else {
        lastAuthResult = LastAuthResult.AuthErrorAcc;
        signResult = SignResult.NotMatchPassword;
      }
    }, 0);
  };
}
