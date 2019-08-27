cmake_minimum_required(VERSION 3.14)
project(q_lang)

set(CMAKE_CXX_STANDARD 17)

file(GLOB_RECURSE SOURCE_FILES "src/*.cpp" "src/*.hpp")

#set(Boost_USE_STATIC_LIBS ON)
#set(Boost_DEBUG ON)
#set(Boost_DETAILED_FAILURE_MSG ON)
set(Boost_COMPILER "-gcc7")
set(Boost_ARCHITECTURE "-x64")
set(BOOST_LIBRARYDIR ${PROJECT_SOURCE_DIR}/boost/lib)
set(BOOST_INCLUDEDIR ${PROJECT_SOURCE_DIR}/boost/include/boost-1_70)
find_package(Boost COMPONENTS system program_options filesystem REQUIRED)
if (Boost_FOUND)
    include_directories(SYSTEM ${Boost_INCLUDE_DIR})
endif ()
include_directories(${PROJECT_SOURCE_DIR}/src)

add_executable(${PROJECT_NAME} ${SOURCE_FILES})
target_link_libraries(${PROJECT_NAME} ${Boost_LIBRARIES})
