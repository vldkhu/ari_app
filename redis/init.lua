-- -- Lua script to initialize Redis with sound files
-- redis.call('LPUSH', 'sound_files', 'sound1.wav')
-- redis.call('LPUSH', 'sound_files', 'sound2.wav')
-- redis.call('LPUSH', 'sound_files', 'sound3.wav')


local soundFiles = {"sound1.wav", "sound2.wav", "sound3.wav"}

for _, file in ipairs(soundFiles) do
    redis.call("LPUSH", "sound_files", file)
end

return "Sound files added to Redis"
