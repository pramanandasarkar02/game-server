import React, { useEffect, useState, useRef, useCallback } from 'react';
import { IoSend } from 'react-icons/io5';
import { FaUserCircle } from 'react-icons/fa';
import { BsDot } from 'react-icons/bs';
import { GoDotFill } from 'react-icons/go';

type Friend = {
  id: string;
  name: string;
  avatar?: string;
  activeStatus: boolean;
};

type Message = {
  id: string;
  senderId: string;
  text: string;
  timestamp: Date;
};

type ChatHistory = {
  [friendId: string]: Message[];
};

const Friends: React.FC = () => {
  const [friends, setFriends] = useState<Friend[]>([]);
  const [activeFriend, setActiveFriend] = useState<Friend | null>(null);
  const [newMessage, setNewMessage] = useState('');
  const [chatHistory, setChatHistory] = useState<ChatHistory>({});
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Sample chat history data
  const sampleChats: ChatHistory = {
    '1': [
      { id: '1', senderId: '1', text: 'Hey there!', timestamp: new Date(Date.now() - 3600000) },
      { id: '2', senderId: 'me', text: 'Hi John! How are you?', timestamp: new Date(Date.now() - 3500000) },
      { id: '3', senderId: '1', text: "I'm good! Just playing some games", timestamp: new Date(Date.now() - 3400000) },
    ],
    '2': [
      { id: '1', senderId: '2', text: 'Are we still meeting tomorrow?', timestamp: new Date(Date.now() - 86400000) },
      { id: '2', senderId: 'me', text: 'Yes, at the usual place', timestamp: new Date(Date.now() - 82800000) },
    ],
    '3': [
      { id: '1', senderId: 'me', text: 'Did you finish that project?', timestamp: new Date(Date.now() - 172800000) },
      { id: '2', senderId: '3', text: 'Almost done! Will send it tonight', timestamp: new Date(Date.now() - 165600000) },
    ],
    '4': [
      { id: '1', senderId: '4', text: 'Check out this cool game I found!', timestamp: new Date(Date.now() - 259200000) },
      { id: '2', senderId: 'me', text: 'Looks awesome!', timestamp: new Date(Date.now() - 258000000) },
      { id: '3', senderId: '4', text: 'We should play together sometime', timestamp: new Date(Date.now() - 257000000) },
    ],
  };

  useEffect(() => {
    const fetchFriends = async () => {
      try {
        const mockFriends: Friend[] = [
          { id: '1', name: 'John Doe', activeStatus: true, avatar: 'https://i.pravatar.cc/150?img=1' },
          { id: '2', name: 'Jane Smith', activeStatus: false, avatar: 'https://i.pravatar.cc/150?img=2' },
          { id: '3', name: 'Alice Johnson', activeStatus: true, avatar: 'https://i.pravatar.cc/150?img=3' },
          { id: '4', name: 'Bob Williams', activeStatus: false, avatar: 'https://i.pravatar.cc/150?img=4' },
        ];

        setFriends(mockFriends);
        setActiveFriend(mockFriends[0]);
        setChatHistory(sampleChats);
      } catch (error) {
        console.error('Failed to fetch friends:', error);
      }
    };

    fetchFriends();
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [activeFriend, chatHistory]);

  const handleFriendClick = useCallback((friend: Friend) => {
    setActiveFriend(friend);
  }, []);

  const handleSendMessage = useCallback(() => {
    if (!newMessage.trim() || !activeFriend) return;

    const message: Message = {
      id: Date.now().toString(),
      senderId: 'me',
      text: newMessage,
      timestamp: new Date(),
    };

    setChatHistory((prev) => ({
      ...prev,
      [activeFriend.id]: [...(prev[activeFriend.id] || []), message],
    }));

    setNewMessage('');
  }, [newMessage, activeFriend]);

  const handleKeyPress = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        handleSendMessage();
      }
    },
    [handleSendMessage]
  );

  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, []);

  const formatTime = useCallback((date: Date) => {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }, []);

  return (
    <div className="flex h-[80vh] max-w-6xl mx-auto rounded-lg overflow-hidden shadow-xl">
      {/* Friends List */}
      <div className="w-[30%] bg-gray-100 border-r border-gray-200 overflow-y-auto">
        <div className="p-5 border-b border-gray-200 bg-white">
          <h2 className="m-0 text-xl font-semibold">Friends</h2>
        </div>
        <ul role="listbox" aria-label="Friends list" className="p-0 m-0 list-none">
          {friends.map((friend) => (
            <li
              key={friend.id}
              onClick={() => handleFriendClick(friend)}
              className={`p-4 cursor-pointer flex items-center border-b border-gray-200 transition-colors duration-200 hover:bg-gray-50 ${activeFriend?.id === friend.id ? 'bg-blue-50' : ''}`}
              role="option"
              aria-selected={activeFriend?.id === friend.id}
              tabIndex={0}
              onKeyPress={(e) => e.key === 'Enter' && handleFriendClick(friend)}
            >
              <div className="relative mr-4">
                {friend.avatar ? (
                  <img
                    src={friend.avatar}
                    alt={`${friend.name}'s avatar`}
                    className="w-10 h-10 rounded-full object-cover"
                  />
                ) : (
                  <FaUserCircle className="text-gray-500 text-4xl" />
                )}
                <BsDot
                  className={`absolute bottom-0 right-0 text-sm bg-white rounded-full ${friend.activeStatus ? 'text-green-500' : 'text-gray-500'}`}
                />
              </div>
              <div className="flex-1">
                <div className="font-medium">{friend.name}</div>
                {chatHistory[friend.id]?.length > 0 && (
                  <div className="text-sm text-gray-500 truncate max-w-[180px]">
                    {chatHistory[friend.id][chatHistory[friend.id].length - 1].text}
                  </div>
                )}
              </div>
            </li>
          ))}
        </ul>
      </div>

      {/* Chat Area */}
      <div className="w-[70%] flex flex-col bg-white">
        {activeFriend ? (
          <>
            {/* Chat Header */}
            <div className="p-4 border-b border-gray-200 flex items-center">
              <div className="relative mr-4">
                {activeFriend.avatar ? (
                  <img
                    src={activeFriend.avatar}
                    alt={`${activeFriend.name}'s avatar`}
                    className="w-10 h-10 rounded-full object-cover"
                  />
                ) : (
                  <FaUserCircle className="text-gray-500 text-4xl" />
                )}
                {/* <BsDot
                  className={`absolute bottom-0 right-0 text-sm bg-gray-700 rounded-full ${activeFriend.activeStatus ? 'text-green-500' : 'text-gray-200'}`}
                /> */}
                {/* <GoDotFill className={`absolute bottom-0 right-0 text-sm bg-gray-500 rounded-full ${activeFriend.activeStatus ? 'text-green-500 bg-green-500' : 'text-gray-500'}`} /> */}
              </div>
              <div>
                <div className="font-medium">{activeFriend.name}</div>
                <div className="text-sm text-gray-500">
                    <GoDotFill className={`inline-block mr-1 text-sm ${activeFriend.activeStatus ? 'text-green-500' : 'text-gray-500'}`} />
                  {activeFriend.activeStatus ? 'Online' : 'Offline'}
                </div>
              </div>
            </div>

            {/* Messages */}
            <div
              className="flex-1 p-5 overflow-y-auto bg-gradient-to-b from-gray-50 to-gray-50"
              role="log"
              aria-label="Chat messages"
            >
              {chatHistory[activeFriend.id]?.map((message) => (
                <div
                  key={message.id}
                  className={`mb-4 flex flex-col ${message.senderId === 'me' ? 'items-end' : 'items-start'}`}
                >
                  <div
                    className={`max-w-[70%] p-3 rounded-2xl shadow-sm ${
                      message.senderId === 'me'
                        ? 'bg-blue-500 text-white rounded-br-none'
                        : 'bg-gray-200 text-black rounded-bl-none'
                    }`}
                  >
                    {message.text}
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    {formatTime(message.timestamp)}
                  </div>
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>

            {/* Message Input */}
            <div className="p-4 border-t border-gray-200 flex items-center">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Type a message..."
                className="flex-1 p-3 rounded-full border border-gray-200 outline-none focus:border-blue-500 text-base"
                aria-label="Type a message"
              />
              <button
                onClick={handleSendMessage}
                disabled={!newMessage.trim()}
                className="ml-3 bg-blue-500 text-white border-none rounded-full w-10 h-10 flex items-center justify-center cursor-pointer transition-colors duration-200 hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
                aria-label="Send message"
              >
                <IoSend />
              </button>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            <p>Select a friend to start chatting</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Friends;