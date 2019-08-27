cmake_minimum_required(VERSION 3.14)
project(q_lang)

set(CMAKE_CXX_STANDARD 17)

file(GLOB SOURCE_FILES "src/*.cpp" "src/*.hpp")

set(Boost_NO_SYSTEM_PATHS ON)
set(Boost_USE_MULTITHREADED ON)
set(Boost_USE_STATIC_LIBS ON)
set(Boost_USE_STATIC_RUNTIME OFF)
set(BOOST_ALL_DYN_LINK OFF)
set(BOOST_ROOT ${PROJECT_SOURCE_DIR}/boost_1_71_0)
find_package(Boost REQUIRED)
if (Boost_FOUND)
    include_directories(SYSTEM ${Boost_INCLUDE_DIR})
endif ()

add_executable(${PROJECT_NAME} ${SOURCE_FILES})
target_link_libraries (${PROJECT_NAME} LINK_PUBLIC Boost::program_options)
