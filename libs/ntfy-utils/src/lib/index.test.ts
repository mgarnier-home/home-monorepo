import { NtfyUtils } from './index.js';

describe('NtfyUtils', () => {
  describe('sendNotification', () => {
    NtfyUtils['config'] = {
      ntfyProtocol: 'http',
      ntfyTopic: 'topic',
      ntfyServer: 'localhost:3000',
    };

    it(`should send a notification`, async () => {
      // Arrange
      const title = 'title';
      const message = 'message';
      const tags = 'tags';

      const fetchMock = jest.fn();
      global.fetch = fetchMock;
      jest.spyOn(Date.prototype, 'toLocaleString').mockReturnValue('01/04/2020 00:00:00');

      // Act
      await NtfyUtils.sendNotification(title, message, tags);

      // Assert
      expect(fetchMock).toHaveBeenCalledWith('http://localhost:3000/topic', {
        body: 'message',
        headers: {
          Tags: 'tags',
          Title: 'title - 01/04/2020 00:00:00',
        },
        method: 'POST',
      });
    });

    it(`shouldnt send a notification if no topic`, async () => {
      // Arrange
      const title = 'title';
      const message = 'message';
      const tags = 'tags';

      NtfyUtils['config'] = {
        ntfyProtocol: 'http',
        ntfyTopic: '',
        ntfyServer: 'localhost:3000',
      };

      const fetchMock = jest.fn();
      global.fetch = fetchMock;

      // Act
      await NtfyUtils.sendNotification(title, message, tags);

      // Assert
      expect(fetchMock).not.toHaveBeenCalled();
    });

    it(`shouldnt send a notification if no server`, async () => {
      // Arrange
      const title = 'title';
      const message = 'message';
      const tags = 'tags';

      NtfyUtils['config'] = {
        ntfyProtocol: 'http',
        ntfyTopic: 'topic',
        ntfyServer: '',
      };

      const fetchMock = jest.fn();
      global.fetch = fetchMock;

      // Act
      await NtfyUtils.sendNotification(title, message, tags);

      // Assert
      expect(fetchMock).not.toHaveBeenCalled();
    });

    it(`shouldnt crash if no config`, async () => {
      // Arrange
      const title = 'title';
      const message = 'message';
      const tags = 'tags';

      NtfyUtils['config'] = undefined as any;

      const fetchMock = jest.fn();
      global.fetch = fetchMock;

      // Act
      await NtfyUtils.sendNotification(title, message, tags);

      // Assert
      expect(fetchMock).not.toHaveBeenCalled();
    });

    it(`shouldnt crash if fetch fails`, async () => {
      // Arrange
      const title = 'title';
      const message = 'message';
      const tags = 'tags';

      NtfyUtils['config'] = {
        ntfyProtocol: 'http',
        ntfyTopic: 'topic',
        ntfyServer: 'localhost:3000',
      };

      const fetchMock = jest.fn().mockRejectedValue('error');
      global.fetch = fetchMock;

      // Act
      await NtfyUtils.sendNotification(title, message, tags);

      // Assert
      expect(fetchMock).toHaveBeenCalled();
    });

    it(`should automatically remove the last / from the server`, async () => {
      // Arrange
      const title = 'title';
      const message = 'message';
      const tags = 'tags';

      NtfyUtils['config'] = {
        ntfyProtocol: 'http',
        ntfyTopic: 'topic',
        ntfyServer: 'localhost:3000/',
      };

      const fetchMock = jest.fn();
      global.fetch = fetchMock;
      jest.spyOn(Date.prototype, 'toLocaleString').mockReturnValue('01/04/2020 00:00:00');

      // Act
      await NtfyUtils.sendNotification(title, message, tags);

      // Assert
      expect(fetchMock).toHaveBeenCalledWith('http://localhost:3000/topic', {
        body: 'message',
        headers: {
          Tags: 'tags',
          Title: 'title - 01/04/2020 00:00:00',
        },
        method: 'POST',
      });
    });

    it(`should automatically remove the first / from the topic`, async () => {
      // Arrange
      const title = 'title';
      const message = 'message';
      const tags = 'tags';

      NtfyUtils['config'] = {
        ntfyProtocol: 'http',
        ntfyTopic: '/topic',
        ntfyServer: 'localhost:3000',
      };

      const fetchMock = jest.fn();
      global.fetch = fetchMock;
      jest.spyOn(Date.prototype, 'toLocaleString').mockReturnValue('01/04/2020 00:00:00');

      // Act
      await NtfyUtils.sendNotification(title, message, tags);

      // Assert
      expect(fetchMock).toHaveBeenCalledWith('http://localhost:3000/topic', {
        body: 'message',
        headers: {
          Tags: 'tags',
          Title: 'title - 01/04/2020 00:00:00',
        },
        method: 'POST',
      });
    });
  });
});
