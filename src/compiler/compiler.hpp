#pragma once

#include <parser/node/master_node.hpp>

namespace ql::compiler {
    class Compiler {
    public:
        void compile(const std::shared_ptr<parser::MasterNode>& program);
    };
}