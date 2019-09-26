#include "compiler.hpp"

#include <iostream>

namespace ql::compiler {
    void Compiler::compile(const std::shared_ptr<parser::MasterNode>& program) {
        parser::AbstractNode::ChildrenRef const& children = program->getChildren();
    }
}
