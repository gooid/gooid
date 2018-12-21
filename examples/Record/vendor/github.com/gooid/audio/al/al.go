// Copyright © 2013-2016 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

// Package al provides golang audio library bindings for OpenAL.
// Official OpenAL documentation can be found online. Prepend "AL_"
// to the function or constant names found in this package.
// Refer to the official OpenAL documentation for more information.
//
// Package al is provided as part of the vu (virtual universe) 3D engine.
package al

// Design Notes:
// These bindings were based on the OpenAL header files found at:
//   http://repo.or.cz/w/openal-soft.git/blob/6dab9d54d1719105e0183f941a2b3dd36e9ba902:/include/AL/al.h
//   http://repo.or.cz/w/openal-soft.git/blob/6dab9d54d1719105e0183f941a2b3dd36e9ba902:/include/AL/alc.h
// Check information available at openal.org.

// //#cgo darwin  LDFLAGS: -framework OpenAL
// //#cgo linux   LDFLAGS: -lopenal -ldl
// //#cgo windows LDFLAGS: -lOpenAL32
//
// #include <stdlib.h>
// #if defined(__APPLE__)
// #include <dlfcn.h>
// #elif defined(_WIN32)
// #define WIN32_LEAN_AND_MEAN 1
// #include <windows.h>
// #else
// #include <dlfcn.h>
// #endif
//
// #ifdef _WIN32
// static HMODULE hmod = NULL;
// #elif !defined __APPLE__
// static void* plib = NULL;
// #endif
// const char* libpath = "libopenal.so";
// static void setLibPath(const char* path) { libpath = path; }
//
// // Helps bind function pointers to c functions.
// static void* bindMethod(const char* name) {
// #ifdef __APPLE__
// 	return dlsym(RTLD_DEFAULT, name);
// #elif _WIN32
// 	if(hmod == NULL) {
// 		hmod = LoadLibraryA("OpenAL32.dll");
// 	}
// 	return GetProcAddress(hmod, (LPCSTR)name);
// #else
// 	if(plib == NULL) {
// 		plib = dlopen(libpath, RTLD_LAZY);
// 	}
// 	return dlsym(plib, name);
// #endif
// }
//
// #if defined(_WIN32)
//  #define AL_APIENTRY __cdecl
//  #define ALC_APIENTRY __cdecl
// #else
//  #define AL_APIENTRY
//  #define ALC_APIENTRY
// #endif
//
// // AL/al.h typedefs
// typedef char ALboolean;
// typedef char ALchar;
// typedef signed char ALbyte;
// typedef unsigned char ALubyte;
// typedef unsigned short ALushort;
// typedef int ALint;
// typedef unsigned int ALuint;
// typedef int ALsizei;
// typedef int ALenum;
// typedef float ALfloat;
// typedef double ALdouble;
// typedef void ALvoid;
//
// #ifndef AL_API
// #define AL_API extern
// #endif
//
// // AL/alc.h typedefs
// typedef struct ALCdevice_struct ALCdevice;
// typedef struct ALCcontext_struct ALCcontext;
// typedef char ALCboolean;
// typedef char ALCchar;
// typedef signed char ALCbyte;
// typedef unsigned char ALCubyte;
// typedef unsigned short ALCushort;
// typedef int ALCint;
// typedef unsigned int ALCuint;
// typedef int ALCsizei;
// typedef int ALCenum;
// typedef void ALCvoid;
//
// #ifndef ALC_API
// #define ALC_API extern
// #endif
//
// // AL/al.h pointers to functions bound to the OS specific library.
// void           (AL_APIENTRY *pfn_alEnable)( ALenum capability );
// void           (AL_APIENTRY *pfn_alDisable)( ALenum capability );
// ALboolean      (AL_APIENTRY *pfn_alIsEnabled)( ALenum capability );
// const ALchar*  (AL_APIENTRY *pfn_alGetString)( ALenum param );
// void           (AL_APIENTRY *pfn_alGetBooleanv)( ALenum param, ALboolean* data );
// void           (AL_APIENTRY *pfn_alGetIntegerv)( ALenum param, ALint* data );
// void           (AL_APIENTRY *pfn_alGetFloatv)( ALenum param, ALfloat* data );
// void           (AL_APIENTRY *pfn_alGetDoublev)( ALenum param, ALdouble* data );
// ALboolean      (AL_APIENTRY *pfn_alGetBoolean)( ALenum param );
// ALint          (AL_APIENTRY *pfn_alGetInteger)( ALenum param );
// ALfloat        (AL_APIENTRY *pfn_alGetFloat)( ALenum param );
// ALdouble       (AL_APIENTRY *pfn_alGetDouble)( ALenum param );
// ALenum         (AL_APIENTRY *pfn_alGetError)( void );
// ALboolean      (AL_APIENTRY *pfn_alIsExtensionPresent)(const ALchar* extname );
// void*          (AL_APIENTRY *pfn_alGetProcAddress)( const ALchar* fname );
// ALenum         (AL_APIENTRY *pfn_alGetEnumValue)( const ALchar* ename );
// void           (AL_APIENTRY *pfn_alListenerf)( ALenum param, ALfloat value );
// void           (AL_APIENTRY *pfn_alListener3f)( ALenum param, ALfloat value1, ALfloat value2, ALfloat value3 );
// void           (AL_APIENTRY *pfn_alListenerfv)( ALenum param, const ALfloat* values );
// void           (AL_APIENTRY *pfn_alListeneri)( ALenum param, ALint value );
// void           (AL_APIENTRY *pfn_alListener3i)( ALenum param, ALint value1, ALint value2, ALint value3 );
// void           (AL_APIENTRY *pfn_alListeneriv)( ALenum param, const ALint* values );
// void           (AL_APIENTRY *pfn_alGetListenerf)( ALenum param, ALfloat* value );
// void           (AL_APIENTRY *pfn_alGetListener3f)( ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3 );
// void           (AL_APIENTRY *pfn_alGetListenerfv)( ALenum param, ALfloat* values );
// void           (AL_APIENTRY *pfn_alGetListeneri)( ALenum param, ALint* value );
// void           (AL_APIENTRY *pfn_alGetListener3i)( ALenum param, ALint *value1, ALint *value2, ALint *value3 );
// void           (AL_APIENTRY *pfn_alGetListeneriv)( ALenum param, ALint* values );
// void           (AL_APIENTRY *pfn_alGenSources)( ALsizei n, ALuint* sources );
// void           (AL_APIENTRY *pfn_alDeleteSources)( ALsizei n, const ALuint* sources );
// ALboolean      (AL_APIENTRY *pfn_alIsSource)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourcef)( ALuint sid, ALenum param, ALfloat value);
// void           (AL_APIENTRY *pfn_alSource3f)( ALuint sid, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3 );
// void           (AL_APIENTRY *pfn_alSourcefv)( ALuint sid, ALenum param, const ALfloat* values );
// void           (AL_APIENTRY *pfn_alSourcei)( ALuint sid, ALenum param, ALint value);
// void           (AL_APIENTRY *pfn_alSource3i)( ALuint sid, ALenum param, ALint value1, ALint value2, ALint value3 );
// void           (AL_APIENTRY *pfn_alSourceiv)( ALuint sid, ALenum param, const ALint* values );
// void           (AL_APIENTRY *pfn_alGetSourcef)( ALuint sid, ALenum param, ALfloat* value );
// void           (AL_APIENTRY *pfn_alGetSource3f)( ALuint sid, ALenum param, ALfloat* value1, ALfloat* value2, ALfloat* value3);
// void           (AL_APIENTRY *pfn_alGetSourcefv)( ALuint sid, ALenum param, ALfloat* values );
// void           (AL_APIENTRY *pfn_alGetSourcei)( ALuint sid, ALenum param, ALint* value );
// void           (AL_APIENTRY *pfn_alGetSource3i)( ALuint sid, ALenum param, ALint* value1, ALint* value2, ALint* value3);
// void           (AL_APIENTRY *pfn_alGetSourceiv)( ALuint sid, ALenum param, ALint* values );
// void           (AL_APIENTRY *pfn_alSourcePlayv)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourceStopv)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourceRewindv)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourcePausev)( ALsizei ns, const ALuint *sids );
// void           (AL_APIENTRY *pfn_alSourcePlay)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourceStop)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourceRewind)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourcePause)( ALuint sid );
// void           (AL_APIENTRY *pfn_alSourceQueueBuffers)(ALuint sid, ALsizei numEntries, const ALuint *bids );
// void           (AL_APIENTRY *pfn_alSourceUnqueueBuffers)(ALuint sid, ALsizei numEntries, ALuint *bids );
// void           (AL_APIENTRY *pfn_alGenBuffers)( ALsizei n, ALuint* buffers );
// void           (AL_APIENTRY *pfn_alDeleteBuffers)( ALsizei n, const ALuint* buffers );
// ALboolean      (AL_APIENTRY *pfn_alIsBuffer)( ALuint bid );
// void           (AL_APIENTRY *pfn_alBufferData)( ALuint bid, ALenum format, const ALvoid* data, ALsizei size, ALsizei freq );
// void           (AL_APIENTRY *pfn_alBufferf)( ALuint bid, ALenum param, ALfloat value);
// void           (AL_APIENTRY *pfn_alBuffer3f)( ALuint bid, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3 );
// void           (AL_APIENTRY *pfn_alBufferfv)( ALuint bid, ALenum param, const ALfloat* values );
// void           (AL_APIENTRY *pfn_alBufferi)( ALuint bid, ALenum param, ALint value);
// void           (AL_APIENTRY *pfn_alBuffer3i)( ALuint bid, ALenum param, ALint value1, ALint value2, ALint value3 );
// void           (AL_APIENTRY *pfn_alBufferiv)( ALuint bid, ALenum param, const ALint* values );
// void           (AL_APIENTRY *pfn_alGetBufferf)( ALuint bid, ALenum param, ALfloat* value );
// void           (AL_APIENTRY *pfn_alGetBuffer3f)( ALuint bid, ALenum param, ALfloat* value1, ALfloat* value2, ALfloat* value3);
// void           (AL_APIENTRY *pfn_alGetBufferfv)( ALuint bid, ALenum param, ALfloat* values );
// void           (AL_APIENTRY *pfn_alGetBufferi)( ALuint bid, ALenum param, ALint* value );
// void           (AL_APIENTRY *pfn_alGetBuffer3i)( ALuint bid, ALenum param, ALint* value1, ALint* value2, ALint* value3);
// void           (AL_APIENTRY *pfn_alGetBufferiv)( ALuint bid, ALenum param, ALint* values );
// void           (AL_APIENTRY *pfn_alDopplerFactor)( ALfloat value );
// void           (AL_APIENTRY *pfn_alDopplerVelocity)( ALfloat value );
// void           (AL_APIENTRY *pfn_alSpeedOfSound)( ALfloat value );
// void           (AL_APIENTRY *pfn_alDistanceModel)( ALenum distanceModel );
//
// // AL/al.h wrappers for the go bindings.
// AL_API void          AL_APIENTRY wrap_alEnable( int capability ) { (*pfn_alEnable)( capability ); }
// AL_API void          AL_APIENTRY wrap_alDisable( int capability ) { (*pfn_alDisable)( capability ); }
// AL_API unsigned int  AL_APIENTRY wrap_alIsEnabled( int capability ) { return (*pfn_alIsEnabled)( capability ); }
// AL_API const char*   AL_APIENTRY wrap_alGetString( int param ) { return (*pfn_alGetString)( param ); }
// AL_API void          AL_APIENTRY wrap_alGetBooleanv( int param, char* data ) { (*pfn_alGetBooleanv)( param, data ); }
// AL_API void          AL_APIENTRY wrap_alGetIntegerv( int param, int* data ) { (*pfn_alGetIntegerv)( param, data ); }
// AL_API void          AL_APIENTRY wrap_alGetFloatv( int param, float* data ) { (*pfn_alGetFloatv)( param, data ); }
// AL_API void          AL_APIENTRY wrap_alGetDoublev( int param, double* data ) { (*pfn_alGetDoublev)( param, data );}
// AL_API ALboolean     AL_APIENTRY wrap_alGetBoolean( int param ) { return (*pfn_alGetBoolean)( param ); }
// AL_API ALint         AL_APIENTRY wrap_alGetInteger( int param ) { return (*pfn_alGetInteger)( param ); }
// AL_API ALfloat       AL_APIENTRY wrap_alGetFloat( int param ) { return (*pfn_alGetFloat)( param ); }
// AL_API ALdouble      AL_APIENTRY wrap_alGetDouble( int param ) { return (*pfn_alGetDouble)( param ); }
// AL_API ALenum        AL_APIENTRY wrap_alGetError( void ) { return (*pfn_alGetError)(); }
// AL_API ALboolean     AL_APIENTRY wrap_alIsExtensionPresent( const char* extname ) { return (*pfn_alIsExtensionPresent)( extname ); }
// AL_API void*         AL_APIENTRY wrap_alGetProcAddress( const char* fname ) { return (*pfn_alGetProcAddress)( fname ); }
// AL_API ALenum        AL_APIENTRY wrap_alGetEnumValue( const char* ename ) { return (*pfn_alGetEnumValue)( ename ); }
// AL_API void          AL_APIENTRY wrap_alListenerf( int param, float value ) { (*pfn_alListenerf)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alListener3f( int param, float value1, float value2, float value3 ) { (*pfn_alListener3f)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alListenerfv( int param, const float* values ) { (*pfn_alListenerfv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alListeneri( int param, int value ) { (*pfn_alListeneri)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alListener3i( int param, int value1, int value2, int value3 ) { (*pfn_alListener3i)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alListeneriv( int param, const int* values ) { (*pfn_alListeneriv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetListenerf( int param, float* value ) { (*pfn_alGetListenerf)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetListener3f( int param, float *value1, float *value2, float *value3 ) { (*pfn_alGetListener3f)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alGetListenerfv( int param, float* values ) { (*pfn_alGetListenerfv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetListeneri( int param, int* value ) { (*pfn_alGetListeneri)( param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetListener3i( int param, int *value1, int *value2, int *value3 ) { (*pfn_alGetListener3i)( param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alGetListeneriv( int param, int* values ) { (*pfn_alGetListeneriv)( param, values ); }
// AL_API void          AL_APIENTRY wrap_alGenSources( int n, unsigned int* sources ) { (*pfn_alGenSources)( n, sources ); }
// AL_API void          AL_APIENTRY wrap_alDeleteSources( int n, const unsigned int* sources ) { (*pfn_alDeleteSources)( n, sources ); }
// AL_API ALboolean     AL_APIENTRY wrap_alIsSource( unsigned int sid ) { return (*pfn_alIsSource)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourcef( unsigned int sid, int param, float value ) { (*pfn_alSourcef)( sid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alSource3f( unsigned int sid, int param, float value1, float value2, float value3 ) { (*pfn_alSource3f)( sid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alSourcefv( unsigned int sid, int param, const float* values ) { (*pfn_alSourcefv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alSourcei( unsigned int sid, int param, int value ) { (*pfn_alSourcei)( sid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alSource3i( unsigned int sid, int param, int value1, int value2, int value3 ) { (*pfn_alSource3i)( sid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alSourceiv( unsigned int sid, int param, const int* values ) { (*pfn_alSourceiv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetSourcef( unsigned int sid, int param, float* value ) { (*pfn_alGetSourcef)( sid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetSource3f( unsigned int sid, int param, float* value1, float* value2, float* value3) { (*pfn_alGetSource3f)( sid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetSourcefv( unsigned int sid, int param, float* values ) { (*pfn_alGetSourcefv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetSourcei( unsigned int sid,  int param, int* value ) { (*pfn_alGetSourcei)( sid,  param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetSource3i( unsigned int sid, int param, int* value1, int* value2, int* value3) { (*pfn_alGetSource3i)( sid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetSourceiv( unsigned int sid,  int param, int* values ) { (*pfn_alGetSourceiv)( sid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alSourcePlayv( int ns, const unsigned int *sids ) { (*pfn_alSourcePlayv)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourceStopv( int ns, const unsigned int *sids ) { (*pfn_alSourceStopv)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourceRewindv( int ns, const unsigned int *sids ) { (*pfn_alSourceRewindv)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourcePausev( int ns, const unsigned int *sids ) { (*pfn_alSourcePausev)( ns, sids ); }
// AL_API void          AL_APIENTRY wrap_alSourcePlay( unsigned int sid ) { (*pfn_alSourcePlay)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourceStop( unsigned int sid ) { (*pfn_alSourceStop)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourceRewind( unsigned int sid ) { (*pfn_alSourceRewind)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourcePause( unsigned int sid ) { (*pfn_alSourcePause)( sid ); }
// AL_API void          AL_APIENTRY wrap_alSourceQueueBuffers( unsigned int sid, int numEntries, const unsigned int *bids ) { (*pfn_alSourceQueueBuffers)( sid, numEntries, bids ); }
// AL_API void          AL_APIENTRY wrap_alSourceUnqueueBuffers( unsigned int sid, int numEntries, unsigned int *bids ) {(*pfn_alSourceUnqueueBuffers)( sid, numEntries, bids ); }
// AL_API void          AL_APIENTRY wrap_alGenBuffers( int n, unsigned int* buffers ) { (*pfn_alGenBuffers)( n, buffers ); }
// AL_API void          AL_APIENTRY wrap_alDeleteBuffers( int n, const unsigned int* buffers ) { (*pfn_alDeleteBuffers)( n, buffers ); }
// AL_API ALboolean     AL_APIENTRY wrap_alIsBuffer( unsigned int bid ) { return (*pfn_alIsBuffer)( bid ); }
// AL_API void          AL_APIENTRY wrap_alBufferData( unsigned int bid, int format, const ALvoid* data, int size, int freq ) { (*pfn_alBufferData)( bid, format, data, size, freq ); }
// AL_API void          AL_APIENTRY wrap_alBufferf( unsigned int bid, int param, float value ) { (*pfn_alBufferf)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alBuffer3f( unsigned int bid, int param, float value1, float value2, float value3 ) { (*pfn_alBuffer3f)( bid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alBufferfv( unsigned int bid, int param, const float* values ) { (*pfn_alBufferfv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alBufferi( unsigned int bid, int param, int value ) { (*pfn_alBufferi)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alBuffer3i( unsigned int bid, int param, int value1, int value2, int value3 ) { (*pfn_alBuffer3i)( bid, param, value1, value2, value3 ); }
// AL_API void          AL_APIENTRY wrap_alBufferiv( unsigned int bid, int param, const int* values ) { (*pfn_alBufferiv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetBufferf( unsigned int bid, int param, float* value ) { (*pfn_alGetBufferf)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetBuffer3f( unsigned int bid, int param, float* value1, float* value2, float* value3) { (*pfn_alGetBuffer3f)( bid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetBufferfv( unsigned int bid, int param, float* values ) { (*pfn_alGetBufferfv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alGetBufferi( unsigned int bid, int param, int* value ) { (*pfn_alGetBufferi)( bid, param, value ); }
// AL_API void          AL_APIENTRY wrap_alGetBuffer3i( unsigned int bid, int param, int* value1, int* value2, int* value3) { (*pfn_alGetBuffer3i)( bid, param, value1, value2, value3); }
// AL_API void          AL_APIENTRY wrap_alGetBufferiv( unsigned int bid, int param, int* values ) { (*pfn_alGetBufferiv)( bid, param, values ); }
// AL_API void          AL_APIENTRY wrap_alDopplerFactor( float value ) { (*pfn_alDopplerFactor)( value ); }
// AL_API void          AL_APIENTRY wrap_alDopplerVelocity( float value ) { (*pfn_alDopplerVelocity)( value ); }
// AL_API void          AL_APIENTRY wrap_alSpeedOfSound( float value ) { (*pfn_alSpeedOfSound)( value ); }
// AL_API void          AL_APIENTRY wrap_alDistanceModel( int distanceModel ) { (*pfn_alDistanceModel)( distanceModel ); }
//
// // AL/alc.h pointers to functions bound to the OS specific library.
// ALCcontext *   (ALC_APIENTRY *pfn_alcCreateContext) (ALCdevice *device, const ALCint *attrlist);
// ALCboolean     (ALC_APIENTRY *pfn_alcMakeContextCurrent)( ALCcontext *context );
// void           (ALC_APIENTRY *pfn_alcProcessContext)( ALCcontext *context );
// void           (ALC_APIENTRY *pfn_alcSuspendContext)( ALCcontext *context );
// void           (ALC_APIENTRY *pfn_alcDestroyContext)( ALCcontext *context );
// ALCcontext *   (ALC_APIENTRY *pfn_alcGetCurrentContext)( void );
// ALCdevice *    (ALC_APIENTRY *pfn_alcGetContextsDevice)( ALCcontext *context );
// ALCdevice *    (ALC_APIENTRY *pfn_alcOpenDevice)( const ALCchar *devicename );
// ALCboolean     (ALC_APIENTRY *pfn_alcCloseDevice)( ALCdevice *device );
// ALCenum        (ALC_APIENTRY *pfn_alcGetError)( ALCdevice *device );
// ALCboolean     (ALC_APIENTRY *pfn_alcIsExtensionPresent)( ALCdevice *device, const ALCchar *extname );
// void *         (ALC_APIENTRY *pfn_alcGetProcAddress)(ALCdevice *device, const ALCchar *funcname );
// ALCenum        (ALC_APIENTRY *pfn_alcGetEnumValue)(ALCdevice *device, const ALCchar *enumname );
// const ALCchar* (ALC_APIENTRY *pfn_alcGetString)( ALCdevice *device, ALCenum param );
// void           (ALC_APIENTRY *pfn_alcGetIntegerv)( ALCdevice *device, ALCenum param, ALCsizei size, ALCint *data );
// ALCdevice *    (ALC_APIENTRY *pfn_alcCaptureOpenDevice)( const ALCchar *devicename, ALCuint frequency, ALCenum format, ALCsizei buffersize );
// ALCboolean     (ALC_APIENTRY *pfn_alcCaptureCloseDevice)( ALCdevice *device );
// void           (ALC_APIENTRY *pfn_alcCaptureStart)( ALCdevice *device );
// void           (ALC_APIENTRY *pfn_alcCaptureStop)( ALCdevice *device );
// void           (ALC_APIENTRY *pfn_alcCaptureSamples)( ALCdevice *device, ALCvoid *buffer, ALCsizei samples );
//
// // AL/alc.h wrappers for the go bindings.
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcCreateContext( uintptr_t device, const int* attrlist ) { return (uintptr_t)(*pfn_alcCreateContext)((ALCdevice *)device, attrlist); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcMakeContextCurrent( uintptr_t context ) { return (*pfn_alcMakeContextCurrent)( (ALCcontext *)context ); }
// ALC_API void         ALC_APIENTRY wrap_alcProcessContext( uintptr_t context ) { (*pfn_alcProcessContext)( (ALCcontext *)context ); }
// ALC_API void         ALC_APIENTRY wrap_alcSuspendContext( uintptr_t context ) { (*pfn_alcSuspendContext)( (ALCcontext *)context ); }
// ALC_API void         ALC_APIENTRY wrap_alcDestroyContext( uintptr_t context ) { (*pfn_alcDestroyContext)( (ALCcontext *)context ); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcGetCurrentContext( void ) { return (uintptr_t)(*pfn_alcGetCurrentContext)(); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcGetContextsDevice( uintptr_t context ) { return (uintptr_t)(*pfn_alcGetContextsDevice)( (ALCcontext *)context ); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcOpenDevice( const char *devicename ) { return (uintptr_t)(*pfn_alcOpenDevice)( devicename ); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcCloseDevice( uintptr_t device ) { return (*pfn_alcCloseDevice)( (ALCdevice *)device ); }
// ALC_API ALCenum      ALC_APIENTRY wrap_alcGetError( uintptr_t device ) { return (*pfn_alcGetError)( (ALCdevice *)device ); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcIsExtensionPresent( uintptr_t device, const char *extname ) { return (*pfn_alcIsExtensionPresent)( (ALCdevice *)device, extname ); }
// ALC_API void  *      ALC_APIENTRY wrap_alcGetProcAddress( uintptr_t device, const char *funcname ) { return (*pfn_alcGetProcAddress)( (ALCdevice *)device, funcname ); }
// ALC_API ALCenum      ALC_APIENTRY wrap_alcGetEnumValue( uintptr_t device, const char *enumname ) { return (*pfn_alcGetEnumValue)( (ALCdevice *)device, enumname ); }
// ALC_API const char * ALC_APIENTRY wrap_alcGetString( uintptr_t device, int param ) { return (*pfn_alcGetString)( (ALCdevice *)device, param ); }
// ALC_API void         ALC_APIENTRY wrap_alcGetIntegerv( uintptr_t device, int param, int size, int *data ) { (*pfn_alcGetIntegerv)( (ALCdevice *)device, param, size, data ); }
// ALC_API uintptr_t    ALC_APIENTRY wrap_alcCaptureOpenDevice( const char *devicename, unsigned int frequency, int format, int buffersize ) { return (uintptr_t)(*pfn_alcCaptureOpenDevice)( devicename, frequency, format, buffersize ); }
// ALC_API ALCboolean   ALC_APIENTRY wrap_alcCaptureCloseDevice( uintptr_t device ) { return (*pfn_alcCaptureCloseDevice)( (ALCdevice *)device ); }
// ALC_API void         ALC_APIENTRY wrap_alcCaptureStart( uintptr_t device ) { (*pfn_alcCaptureStart)( (ALCdevice *)device ); }
// ALC_API void         ALC_APIENTRY wrap_alcCaptureStop( uintptr_t device ) { (*pfn_alcCaptureStop)( (ALCdevice *)device ); }
// ALC_API void         ALC_APIENTRY wrap_alcCaptureSamples( uintptr_t device, ALCvoid *buffer, int samples ) { (*pfn_alcCaptureSamples)( (ALCdevice *)device, buffer, samples ); }
//
// void al_init() {
//    // AL/al.h
//    pfn_alEnable                  = bindMethod("alEnable");
//    pfn_alDisable                 = bindMethod("alDisable");
//    pfn_alIsEnabled               = bindMethod("alIsEnabled");
//    pfn_alGetString               = bindMethod("alGetString");
//    pfn_alGetBooleanv             = bindMethod("alGetBooleanv");
//    pfn_alGetIntegerv             = bindMethod("alGetIntegerv");
//    pfn_alGetFloatv               = bindMethod("alGetFloatv");
//    pfn_alGetDoublev              = bindMethod("alGetDoublev");
//    pfn_alGetBoolean              = bindMethod("alGetBoolean");
//    pfn_alGetInteger              = bindMethod("alGetInteger");
//    pfn_alGetFloat                = bindMethod("alGetFloat");
//    pfn_alGetDouble               = bindMethod("alGetDouble");
//    pfn_alGetError                = bindMethod("alGetError");
//    pfn_alIsExtensionPresent      = bindMethod("alIsExtensionPresent");
//    pfn_alGetProcAddress          = bindMethod("alGetProcAddress");
//    pfn_alGetEnumValue            = bindMethod("alGetEnumValue");
//    pfn_alListenerf               = bindMethod("alListenerf");
//    pfn_alListener3f              = bindMethod("alListener3f");
//    pfn_alListenerfv              = bindMethod("alListenerfv");
//    pfn_alListeneri               = bindMethod("alListeneri");
//    pfn_alListener3i              = bindMethod("alListener3i");
//    pfn_alListeneriv              = bindMethod("alListeneriv");
//    pfn_alGetListenerf            = bindMethod("alGetListenerf");
//    pfn_alGetListener3f           = bindMethod("alGetListener3f");
//    pfn_alGetListenerfv           = bindMethod("alGetListenerfv");
//    pfn_alGetListeneri            = bindMethod("alGetListeneri");
//    pfn_alGetListener3i           = bindMethod("alGetListener3i");
//    pfn_alGetListeneriv           = bindMethod("alGetListeneriv");
//    pfn_alGenSources              = bindMethod("alGenSources");
//    pfn_alDeleteSources           = bindMethod("alDeleteSources");
//    pfn_alIsSource                = bindMethod("alIsSource");
//    pfn_alSourcef                 = bindMethod("alSourcef");
//    pfn_alSource3f                = bindMethod("alSource3f");
//    pfn_alSourcefv                = bindMethod("alSourcefv");
//    pfn_alSourcei                 = bindMethod("alSourcei");
//    pfn_alSource3i                = bindMethod("alSource3i");
//    pfn_alSourceiv                = bindMethod("alSourceiv");
//    pfn_alGetSourcef              = bindMethod("alGetSourcef");
//    pfn_alGetSource3f             = bindMethod("alGetSource3f");
//    pfn_alGetSourcefv             = bindMethod("alGetSourcefv");
//    pfn_alGetSourcei              = bindMethod("alGetSourcei");
//    pfn_alGetSource3i             = bindMethod("alGetSource3i");
//    pfn_alGetSourceiv             = bindMethod("alGetSourceiv");
//    pfn_alSourcePlayv             = bindMethod("alSourcePlayv");
//    pfn_alSourceStopv             = bindMethod("alSourceStopv");
//    pfn_alSourceRewindv           = bindMethod("alSourceRewindv");
//    pfn_alSourcePausev            = bindMethod("alSourcePausev");
//    pfn_alSourcePlay              = bindMethod("alSourcePlay");
//    pfn_alSourceStop              = bindMethod("alSourceStop");
//    pfn_alSourceRewind            = bindMethod("alSourceRewind");
//    pfn_alSourcePause             = bindMethod("alSourcePause");
//    pfn_alSourceQueueBuffers      = bindMethod("alSourceQueueBuffers");
//    pfn_alSourceUnqueueBuffers    = bindMethod("alSourceUnqueueBuffers");
//    pfn_alGenBuffers              = bindMethod("alGenBuffers");
//    pfn_alDeleteBuffers           = bindMethod("alDeleteBuffers");
//    pfn_alIsBuffer                = bindMethod("alIsBuffer");
//    pfn_alBufferData              = bindMethod("alBufferData");
//    pfn_alBufferf                 = bindMethod("alBufferf");
//    pfn_alBuffer3f                = bindMethod("alBuffer3f");
//    pfn_alBufferfv                = bindMethod("alBufferfv");
//    pfn_alBufferi                 = bindMethod("alBufferi");
//    pfn_alBuffer3i                = bindMethod("alBuffer3i");
//    pfn_alBufferiv                = bindMethod("alBufferiv");
//    pfn_alGetBufferf              = bindMethod("alGetBufferf");
//    pfn_alGetBuffer3f             = bindMethod("alGetBuffer3f");
//    pfn_alGetBufferfv             = bindMethod("alGetBufferfv");
//    pfn_alGetBufferi              = bindMethod("alGetBufferi");
//    pfn_alGetBuffer3i             = bindMethod("alGetBuffer3i");
//    pfn_alGetBufferiv             = bindMethod("alGetBufferiv");
//    pfn_alDopplerFactor           = bindMethod("alDopplerFactor");
//    pfn_alDopplerVelocity         = bindMethod("alDopplerVelocity");
//    pfn_alSpeedOfSound            = bindMethod("alSpeedOfSound");
//    pfn_alDistanceModel           = bindMethod("alDistanceModel");
//
//    // AL/alc.h
//    pfn_alcCreateContext          = bindMethod("alcCreateContext");
//    pfn_alcMakeContextCurrent     = bindMethod("alcMakeContextCurrent");
//    pfn_alcProcessContext         = bindMethod("alcProcessContext");
//    pfn_alcSuspendContext         = bindMethod("alcSuspendContext");
//    pfn_alcDestroyContext         = bindMethod("alcDestroyContext");
//    pfn_alcGetCurrentContext      = bindMethod("alcGetCurrentContext");
//    pfn_alcGetContextsDevice      = bindMethod("alcGetContextsDevice");
//    pfn_alcOpenDevice             = bindMethod("alcOpenDevice");
//    pfn_alcCloseDevice            = bindMethod("alcCloseDevice");
//    pfn_alcGetError               = bindMethod("alcGetError");
//    pfn_alcIsExtensionPresent     = bindMethod("alcIsExtensionPresent");
//    pfn_alcGetProcAddress         = bindMethod("alcGetProcAddress");
//    pfn_alcGetEnumValue           = bindMethod("alcGetEnumValue");
//    pfn_alcGetString              = bindMethod("alcGetString");
//    pfn_alcGetIntegerv            = bindMethod("alcGetIntegerv");
//    pfn_alcCaptureOpenDevice      = bindMethod("alcCaptureOpenDevice");
//    pfn_alcCaptureCloseDevice     = bindMethod("alcCaptureCloseDevice");
//    pfn_alcCaptureStart           = bindMethod("alcCaptureStart");
//    pfn_alcCaptureStop            = bindMethod("alcCaptureStop");
//    pfn_alcCaptureSamples         = bindMethod("alcCaptureSamples");
// }
//
import "C"
import (
	"errors"
	"strings"
	"unsafe"
)
import "fmt"

// AL/al.h constants (with AL_ removed). Refer to the original header for constant documentation.
const (
	FALSE                     = 0
	TRUE                      = 1
	NONE                      = 0
	NO_ERROR                  = 0
	SOURCE_RELATIVE           = 0x202
	CONE_INNER_ANGLE          = 0x1001
	CONE_OUTER_ANGLE          = 0x1002
	PITCH                     = 0x1003
	POSITION                  = 0x1004
	DIRECTION                 = 0x1005
	VELOCITY                  = 0x1006
	LOOPING                   = 0x1007
	BUFFER                    = 0x1009
	GAIN                      = 0x100A
	MIN_GAIN                  = 0x100D
	MAX_GAIN                  = 0x100E
	ORIENTATION               = 0x100F
	SOURCE_STATE              = 0x1010
	INITIAL                   = 0x1011
	PLAYING                   = 0x1012
	PAUSED                    = 0x1013
	STOPPED                   = 0x1014
	BUFFERS_QUEUED            = 0x1015
	BUFFERS_PROCESSED         = 0x1016
	SEC_OFFSET                = 0x1024
	SAMPLE_OFFSET             = 0x1025
	BYTE_OFFSET               = 0x1026
	SOURCE_TYPE               = 0x1027
	STATIC                    = 0x1028
	STREAMING                 = 0x1029
	UNDETERMINED              = 0x1030
	FORMAT_MONO8              = 0x1100
	FORMAT_MONO16             = 0x1101
	FORMAT_STEREO8            = 0x1102
	FORMAT_STEREO16           = 0x1103
	REFERENCE_DISTANCE        = 0x1020
	ROLLOFF_FACTOR            = 0x1021
	CONE_OUTER_GAIN           = 0x1022
	MAX_DISTANCE              = 0x1023
	FREQUENCY                 = 0x2001
	BITS                      = 0x2002
	CHANNELS                  = 0x2003
	SIZE                      = 0x2004
	UNUSED                    = 0x2010
	PENDING                   = 0x2011
	PROCESSED                 = 0x2012
	INVALID_NAME              = 0xA001
	INVALID_ENUM              = 0xA002
	INVALID_VALUE             = 0xA003
	INVALID_OPERATION         = 0xA004
	OUT_OF_MEMORY             = 0xA005
	VENDOR                    = 0xB001
	VERSION                   = 0xB002
	RENDERER                  = 0xB003
	EXTENSIONS                = 0xB004
	DOPPLER_FACTOR            = 0xC000
	DOPPLER_VELOCITY          = 0xC001
	SPEED_OF_SOUND            = 0xC003
	DISTANCE_MODEL            = 0xD000
	INVERSE_DISTANCE          = 0xD001
	INVERSE_DISTANCE_CLAMPED  = 0xD002
	LINEAR_DISTANCE           = 0xD003
	LINEAR_DISTANCE_CLAMPED   = 0xD004
	EXPONENT_DISTANCE         = 0xD005
	EXPONENT_DISTANCE_CLAMPED = 0xD006
)

// AL/alc.h constants (with AL removed). Refer to the original header for constant documentation.
const (
	C_FALSE                            = 0
	C_TRUE                             = 1
	C_NO_ERROR                         = 0
	C_FREQUENCY                        = 0x1007
	C_REFRESH                          = 0x1008
	C_SYNC                             = 0x1009
	C_MONO_SOURCES                     = 0x1010
	C_STEREO_SOURCES                   = 0x1011
	C_INVALID_DEVICE                   = 0xA001
	C_INVALID_CONTEXT                  = 0xA002
	C_INVALID_ENUM                     = 0xA003
	C_INVALID_VALUE                    = 0xA004
	C_OUT_OF_MEMORY                    = 0xA005
	C_DEFAULT_DEVICE_SPECIFIER         = 0x1004
	C_DEVICE_SPECIFIER                 = 0x1005
	C_EXTENSIONS                       = 0x1006
	C_MAJOR_VERSION                    = 0x1000
	C_MINOR_VERSION                    = 0x1001
	C_ATTRIBUTES_SIZE                  = 0x1002
	C_ALL_ATTRIBUTES                   = 0x1003
	C_CAPTURE_DEVICE_SPECIFIER         = 0x310
	C_CAPTURE_DEFAULT_DEVICE_SPECIFIER = 0x311
	C_CAPTURE_SAMPLES                  = 0x312
)

var bInit = false
var bindError error

// bind the methods to the function pointers
func Init() error {
	if !bInit {
		C.al_init()
		bInit = true
		ss := BindingReport()
		estr := ""
		for _, s := range ss {
			if strings.Index(s, "[ ]") >= 0 {
				estr += s + ";"
			}
		}
		if estr != "" {
			bindError = errors.New(estr)
		}
	}
	return bindError
}

func InitPath(path string) error {
	if !bInit {
		if path != "" {
			cstr := C.CString(path)
			defer C.free(unsafe.Pointer(cstr))
			C.setLibPath(cstr)
		}
		Init()
	}
	return bindError
}

// Show which function pointers are bound [+] or not bound [-].
// Expected to be used as a sanity check to see if the OpenAL libraries exist.
func BindingReport() (report []string) {
	report = []string{}

	// AL/al.h
	report = append(report, "AL")
	report = append(report, isBound(unsafe.Pointer(C.pfn_alEnable), "alEnable"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDisable), "alDisable"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsEnabled), "alIsEnabled"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetString), "alGetString"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBooleanv), "alGetBooleanv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetIntegerv), "alGetIntegerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetFloatv), "alGetFloatv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetDoublev), "alGetDoublev"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBoolean), "alGetBoolean"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetInteger), "alGetInteger"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetFloat), "alGetFloat"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetDouble), "alGetDouble"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetError), "alGetError"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsExtensionPresent), "alIsExtensionPresent"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetProcAddress), "alGetProcAddress"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetEnumValue), "alGetEnumValue"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListenerf), "alListenerf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListener3f), "alListener3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListenerfv), "alListenerfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListeneri), "alListeneri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListener3i), "alListener3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alListeneriv), "alListeneriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListenerf), "alGetListenerf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListener3f), "alGetListener3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListenerfv), "alGetListenerfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListeneri), "alGetListeneri"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListener3i), "alGetListener3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetListeneriv), "alGetListeneriv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGenSources), "alGenSources"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDeleteSources), "alDeleteSources"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsSource), "alIsSource"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcef), "alSourcef"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSource3f), "alSource3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcefv), "alSourcefv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcei), "alSourcei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSource3i), "alSource3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceiv), "alSourceiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourcef), "alGetSourcef"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSource3f), "alGetSource3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourcefv), "alGetSourcefv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourcei), "alGetSourcei"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSource3i), "alGetSource3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetSourceiv), "alGetSourceiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePlayv), "alSourcePlayv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceStopv), "alSourceStopv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceRewindv), "alSourceRewindv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePausev), "alSourcePausev"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePlay), "alSourcePlay"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceStop), "alSourceStop"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceRewind), "alSourceRewind"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourcePause), "alSourcePause"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceQueueBuffers), "alSourceQueueBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSourceUnqueueBuffers), "alSourceUnqueueBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGenBuffers), "alGenBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDeleteBuffers), "alDeleteBuffers"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alIsBuffer), "alIsBuffer"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferData), "alBufferData"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferf), "alBufferf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBuffer3f), "alBuffer3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferfv), "alBufferfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferi), "alBufferi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBuffer3i), "alBuffer3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alBufferiv), "alBufferiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferf), "alGetBufferf"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBuffer3f), "alGetBuffer3f"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferfv), "alGetBufferfv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferi), "alGetBufferi"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBuffer3i), "alGetBuffer3i"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alGetBufferiv), "alGetBufferiv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDopplerFactor), "alDopplerFactor"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDopplerVelocity), "alDopplerVelocity"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alSpeedOfSound), "alSpeedOfSound"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alDistanceModel), "alDistanceModel"))

	// AL/alc.h
	report = append(report, "ALC")
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCreateContext), "alcCreateContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcMakeContextCurrent), "alcMakeContextCurrent"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcProcessContext), "alcProcessContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcSuspendContext), "alcSuspendContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcDestroyContext), "alcDestroyContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetCurrentContext), "alcGetCurrentContext"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetContextsDevice), "alcGetContextsDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcOpenDevice), "alcOpenDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCloseDevice), "alcCloseDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetError), "alcGetError"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcIsExtensionPresent), "alcIsExtensionPresent"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetProcAddress), "alcGetProcAddress"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetEnumValue), "alcGetEnumValue"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetString), "alcGetString"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcGetIntegerv), "alcGetIntegerv"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureOpenDevice), "alcCaptureOpenDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureCloseDevice), "alcCaptureCloseDevice"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureStart), "alcCaptureStart"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureStop), "alcCaptureStop"))
	report = append(report, isBound(unsafe.Pointer(C.pfn_alcCaptureSamples), "alcCaptureSamples"))
	return
}

// isBound returns a string that indicates if the given function
// pointer is bound.
func isBound(pfn unsafe.Pointer, fn string) string {
	inc := " "
	if pfn != nil {
		inc = "+"
	}
	return fmt.Sprintf("   [%s] %s", inc, fn)
}

// al.h
// Enable enables a capability.
func Enable(capability int32) {
	C.wrap_alEnable(C.ALenum(capability))
}

// Disable disables a capability.
func Disable(capability int32) {
	C.wrap_alDisable(C.ALenum(capability))
}

// Enabled returns true if the specified capability is enabled.
func Enabled(capability int32) bool {
	return C.wrap_alIsEnabled(C.ALenum(capability)) == 1
}

// Vector represents an vector in a Cartesian coordinate system.
type Vector [3]float32

// Orientation represents the angular position of an object in a
// right-handed Cartesian coordinate system.
// A cross product between the forward and up vector returns a vector
// that points to the right.
type Orientation struct {
	// Forward vector is the direction that the object is looking at.
	Forward Vector
	// Up vector represents the rotation of the object.
	Up Vector
}

func orientationFromSlice(v []float32) Orientation {
	return Orientation{
		Forward: Vector{v[0], v[1], v[2]},
		Up:      Vector{v[3], v[4], v[5]},
	}
}

func (v Orientation) slice() []float32 {
	return []float32{v.Forward[0], v.Forward[1], v.Forward[2], v.Up[0], v.Up[1], v.Up[2]}
}

func geti(param int) int32 {
	return int32(C.wrap_alGetInteger(C.ALenum(param)))
}

func getf(param int) float32 {
	return float32(C.wrap_alGetFloat(C.ALenum(param)))
}

func getString(param int) string {
	v := C.wrap_alGetString(C.ALenum(param))
	return C.GoString((*C.char)(v))
}

// DistanceModel returns the distance model.
func DistanceModel() int32 {
	return geti(paramDistanceModel)
}

// SetDistanceModel sets the distance model.
func SetDistanceModel(v int32) {
	C.wrap_alDistanceModel(C.ALenum(v))
}

// DopplerFactor returns the doppler factor.
func DopplerFactor() float32 {
	return getf(paramDopplerFactor)
}

// SetDopplerFactor sets the doppler factor.
func SetDopplerFactor(v float32) {
	C.wrap_alDopplerFactor(C.ALfloat(v))
}

// DopplerVelocity returns the doppler velocity.
func DopplerVelocity() float32 {
	return getf(paramDopplerVelocity)
}

// SetDopplerVelocity sets the doppler velocity.
func SetDopplerVelocity(v float32) {
	C.wrap_alDopplerVelocity(C.ALfloat(v))
}

// SpeedOfSound is the speed of sound in meters per second (m/s).
func SpeedOfSound() float32 {
	return getf(paramSpeedOfSound)
}

// SetSpeedOfSound sets the speed of sound, its unit should be meters per second (m/s).
func SetSpeedOfSound(v float32) {
	C.wrap_alSpeedOfSound(C.ALfloat(v))
}

// Vendor returns the vendor.
func Vendor() string {
	return getString(paramVendor)
}

// Version returns the version string.
func Version() string {
	return getString(paramVersion)
}

// Renderer returns the renderer information.
func Renderer() string {
	return getString(paramRenderer)
}

// Extensions returns the enabled extensions.
func Extensions() string {
	return getString(paramExtensions)
}

// Error returns the most recently generated error.
func Error() int32 {
	return int32(C.wrap_alGetError())
}

// Source represents an individual sound source in 3D-space.
// They take PCM data, apply modifications and then submit them to
// be mixed according to their spatial location.
type Source uint32

// GenSources generates n new sources. These sources should be deleted
// once they are not in use.
func GenSources(n int) []Source {
	s := make([]Source, n)
	C.wrap_alGenSources(C.ALsizei(n), (*C.ALuint)(unsafe.Pointer(&s[0])))
	return s
}

// PlaySources plays the sources.
func PlaySources(source ...Source) {
	C.wrap_alSourcePlayv(C.ALsizei(len(source)), (*C.ALuint)(unsafe.Pointer(&source[0])))
}

// PauseSources pauses the sources.
func PauseSources(source ...Source) {
	C.wrap_alSourcePausev(C.ALsizei(len(source)), (*C.ALuint)(unsafe.Pointer(&source[0])))
}

// StopSources stops the sources.
func StopSources(source ...Source) {
	C.wrap_alSourceStopv(C.ALsizei(len(source)), (*C.ALuint)(unsafe.Pointer(&source[0])))
}

// RewindSources rewinds the sources to their beginning positions.
func RewindSources(source ...Source) {
	C.wrap_alSourceRewindv(C.ALsizei(len(source)), (*C.ALuint)(unsafe.Pointer(&source[0])))
}

// DeleteSources deletes the sources.
func DeleteSources(source ...Source) {
	C.wrap_alDeleteSources(C.ALsizei(len(source)), (*C.ALuint)(unsafe.Pointer(&source[0])))
}

// Gain returns the source gain.
func (s Source) Gain() float32 {
	return getSourcef(s, paramGain)
}

// SetGain sets the source gain.
func (s Source) SetGain(v float32) {
	setSourcef(s, paramGain, v)
}

// MinGain returns the source's minimum gain setting.
func (s Source) MinGain() float32 {
	return getSourcef(s, paramMinGain)
}

// SetMinGain sets the source's minimum gain setting.
func (s Source) SetMinGain(v float32) {
	setSourcef(s, paramMinGain, v)
}

// MaxGain returns the source's maximum gain setting.
func (s Source) MaxGain() float32 {
	return getSourcef(s, paramMaxGain)
}

// SetMaxGain sets the source's maximum gain setting.
func (s Source) SetMaxGain(v float32) {
	setSourcef(s, paramMaxGain, v)
}

// Position returns the position of the source.
func (s Source) Position() Vector {
	v := Vector{}
	getSourcefv(s, paramPosition, v[:])
	return v
}

// SetPosition sets the position of the source.
func (s Source) SetPosition(v Vector) {
	setSourcefv(s, paramPosition, v[:])
}

// Velocity returns the source's velocity.
func (s Source) Velocity() Vector {
	v := Vector{}
	getSourcefv(s, paramVelocity, v[:])
	return v
}

// SetVelocity sets the source's velocity.
func (s Source) SetVelocity(v Vector) {
	setSourcefv(s, paramVelocity, v[:])
}

// Orientation returns the orientation of the source.
func (s Source) Orientation() Orientation {
	v := make([]float32, 6)
	getSourcefv(s, paramOrientation, v)
	return orientationFromSlice(v)
}

// SetOrientation sets the orientation of the source.
func (s Source) SetOrientation(o Orientation) {
	setSourcefv(s, paramOrientation, o.slice())
}

// State returns the playing state of the source.
func (s Source) State() int32 {
	return getSourcei(s, paramSourceState)
}

// BuffersQueued returns the number of the queued buffers.
func (s Source) BuffersQueued() int32 {
	return getSourcei(s, paramBuffersQueued)
}

// BuffersProcessed returns the number of the processed buffers.
func (s Source) BuffersProcessed() int32 {
	return getSourcei(s, paramBuffersProcessed)
}

// OffsetSeconds returns the current playback position of the source in seconds.
func (s Source) OffsetSeconds() int32 {
	return getSourcei(s, paramSecOffset)
}

// OffsetSample returns the sample offset of the current playback position.
func (s Source) OffsetSample() int32 {
	return getSourcei(s, paramSampleOffset)
}

// OffsetByte returns the byte offset of the current playback position.
func (s Source) OffsetByte() int32 {
	return getSourcei(s, paramByteOffset)
}

func getSourcei(s Source, param int) int32 {
	var v C.ALint
	C.wrap_alGetSourcei(C.ALuint(s), C.ALenum(param), &v)
	return int32(v)
}

func getSourcef(s Source, param int) float32 {
	var v C.ALfloat
	C.wrap_alGetSourcef(C.ALuint(s), C.ALenum(param), &v)
	return float32(v)
}

func getSourcefv(s Source, param int, v []float32) {
	C.wrap_alGetSourcefv(C.ALuint(s), C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&v[0])))
}

func setSourcei(s Source, param int, v int32) {
	C.wrap_alSourcei(C.ALuint(s), C.ALenum(param), C.ALint(v))
}

func setSourcef(s Source, param int, v float32) {
	C.wrap_alSourcef(C.ALuint(s), C.ALenum(param), C.ALfloat(v))
}

func setSourcefv(s Source, param int, v []float32) {
	C.wrap_alSourcefv(C.ALuint(s), C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&v[0])))
}

// QueueBuffers adds the buffers to the buffer queue.
func (s Source) QueueBuffers(buffers []Buffer) {
	C.wrap_alSourceQueueBuffers(C.ALuint(s), C.ALsizei(len(buffers)), (*C.ALuint)(unsafe.Pointer(&buffers[0])))
}

// UnqueueBuffers removes the specified buffers from the buffer queue.
func (s Source) UnqueueBuffers(buffers []Buffer) {
	C.wrap_alSourceUnqueueBuffers(C.ALuint(s), C.ALsizei(len(buffers)), (*C.ALuint)(unsafe.Pointer(&buffers[0])))
}

// ListenerGain returns the total gain applied to the final mix.
func ListenerGain() float32 {
	return getListenerf(paramGain)
}

// ListenerPosition returns the position of the listener.
func ListenerPosition() Vector {
	v := Vector{}
	getListenerfv(paramPosition, v[:])
	return v
}

// ListenerVelocity returns the velocity of the listener.
func ListenerVelocity() Vector {
	v := Vector{}
	getListenerfv(paramVelocity, v[:])
	return v
}

// ListenerOrientation returns the orientation of the listener.
func ListenerOrientation() Orientation {
	v := make([]float32, 6)
	getListenerfv(paramOrientation, v)
	return orientationFromSlice(v)
}

// SetListenerGain sets the total gain that will be applied to the final mix.
func SetListenerGain(v float32) {
	setListenerf(paramGain, v)
}

// SetListenerPosition sets the position of the listener.
func SetListenerPosition(v Vector) {
	setListenerfv(paramPosition, v[:])
}

// SetListenerVelocity sets the velocity of the listener.
func SetListenerVelocity(v Vector) {
	setListenerfv(paramVelocity, v[:])
}

// SetListenerOrientation sets the orientation of the listener.
func SetListenerOrientation(v Orientation) {
	setListenerfv(paramOrientation, v.slice())
}

func getListenerf(param int) float32 {
	var v C.ALfloat
	C.wrap_alGetListenerf(C.ALenum(param), &v)
	return float32(v)
}

func getListenerfv(param int, v []float32) {
	C.wrap_alGetListenerfv(C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&v[0])))
}

func setListenerf(param int, v float32) {
	C.wrap_alListenerf(C.ALenum(param), C.ALfloat(v))
}

func setListenerfv(param int, v []float32) {
	C.wrap_alListenerfv(C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&v[0])))
}

// A buffer represents a chunk of PCM audio data that could be buffered to an audio
// source. A single buffer could be shared between multiple sources.
type Buffer uint32

// GenBuffers generates n new buffers. The generated buffers should be deleted
// once they are no longer in use.
func GenBuffers(n int) []Buffer {
	s := make([]Buffer, n)
	C.wrap_alGenBuffers(C.ALsizei(n), (*C.ALuint)(unsafe.Pointer(&s[0])))
	return s
}

// DeleteBuffers deletes the buffers.
func DeleteBuffers(buffers []Buffer) {
	C.wrap_alDeleteBuffers(C.ALsizei(len(buffers)), (*C.ALuint)(unsafe.Pointer(&buffers[0])))
}

func getBufferi(b Buffer, param int) int32 {
	var v C.ALint
	C.wrap_alGetBufferi(C.ALuint(b), C.ALenum(param), &v)
	return int32(v)
}

// Frequency returns the frequency of the buffer data in Hertz (Hz).
func (b Buffer) Frequency() int32 {
	return getBufferi(b, paramFreq)
}

// Bits return the number of bits used to represent a sample.
func (b Buffer) Bits() int32 {
	return getBufferi(b, paramBits)
}

// Channels return the number of the audio channels.
func (b Buffer) Channels() int32 {
	return getBufferi(b, paramChannels)
}

// Size returns the size of the data.
func (b Buffer) Size() int32 {
	return getBufferi(b, paramSize)
}

// BufferData buffers PCM data to the current buffer.
func (b Buffer) BufferData(format uint32, data []byte, freq int32) {
	C.wrap_alBufferData(C.ALuint(b), C.ALenum(format), unsafe.Pointer(&data[0]), C.ALsizei(len(data)), C.ALsizei(freq))
}

// Valid returns true if the buffer exists and is valid.
func (b Buffer) Valid() bool {
	return C.wrap_alIsBuffer(C.ALuint(b)) == 1
}

// alc.h
// Device represents an audio device.
type Device struct {
	d C.uintptr_t
}

// Error returns the last known error from the current device.
func (d *Device) Error() int32 {
	return int32(C.wrap_alcGetError(d.d))
}

// Context represents a context created in the OpenAL layer. A valid current
// context is required to run OpenAL functions.
// The returned context will be available process-wide if it's made the
// current by calling MakeContextCurrent.
type Context struct {
	c C.uintptr_t
}

// Open opens a new device in the OpenAL layer.
func Open(name string) *Device {
	if Init() != nil {
		return nil
	}
	var n *C.char
	if name != "" {
		n = C.CString(name)
		defer C.free(unsafe.Pointer(n))
	}
	d := C.wrap_alcOpenDevice((*C.ALCchar)(unsafe.Pointer(n)))
	if d == 0 {
		return nil
	}
	return &Device{d: d}
}

// Close closes the device.
func (d *Device) Close() bool {
	return C.wrap_alcCloseDevice(d.d) != 0
}

// CreateContext creates a new context.
func (d *Device) CreateContext(attrs []int32) *Context {
	// TODO(jbd): Handle attributes.
	c := C.wrap_alcCreateContext(d.d, nil)
	if c == 0 {
		return nil
	}
	return &Context{c: c}
}

// MakeContextCurrent makes a context current process-wide.
func (c *Context) MakeCurrent() bool {
	return C.wrap_alcMakeContextCurrent(c.c) != 0
}

func (context Context) Process() {
	C.wrap_alcProcessContext(context.c)
}

func (context Context) Suspend() {
	C.wrap_alcSuspendContext(context.c)
}

func (context Context) Destroy() {
	C.wrap_alcDestroyContext(context.c)
}

func GetCurrentContext() *Context {
	c := C.wrap_alcGetCurrentContext()
	if c == 0 {
		return nil
	}
	return &Context{c: c}
}

func (context Context) GetDevice() *Device {
	d := C.wrap_alcGetContextsDevice(context.c)
	if d == 0 {
		return nil
	}
	return &Device{d: d}
}

//
func (d *Device) GetEnumValue(enumname string) int {
	cenumname := C.CString(enumname)
	defer C.free(unsafe.Pointer(cenumname))
	return int(C.wrap_alcGetEnumValue(d.d, cenumname))
}

func (d *Device) GetString(param int) string {
	return C.GoString(C.wrap_alcGetString(d.d, C.ALCenum(param)))
}

func (d *Device) GetIntegerv(param int, size int64, values *int32) {
	C.wrap_alcGetIntegerv(d.d, C.ALCenum(param), C.ALCsizei(size), (*C.ALCint)(unsafe.Pointer(values)))
}

// Capture device
type CaptureDevice struct {
	Device
}

// Open opens a new capture device in the OpenAL layer.
func CaptureOpen(name string, frequency uint, format int, buffersize int64) *CaptureDevice {
	if Init() != nil {
		return nil
	}
	var n *C.char
	if name != "" {
		n = C.CString(name)
		defer C.free(unsafe.Pointer(n))
	}
	d := C.wrap_alcCaptureOpenDevice((*C.ALCchar)(unsafe.Pointer(n)),
		C.ALCuint(frequency), C.ALCenum(format), C.ALCsizei(buffersize))
	if d == 0 {
		return nil
	}
	return &CaptureDevice{Device: Device{d: d}}
}

// Close closes the capture device.
func (d *CaptureDevice) Close() bool {
	return C.wrap_alcCaptureCloseDevice(d.d) != 0
}

func (dev *CaptureDevice) Start() {
	C.wrap_alcCaptureStart(dev.d)
}

func (dev *CaptureDevice) Stop() {
	C.wrap_alcCaptureStop(dev.d)
}

func (dev *CaptureDevice) Samples(bs []byte, samples int64) {
	C.wrap_alcCaptureSamples(dev.d, unsafe.Pointer(&bs[0]), C.ALCsizei(samples))
}
