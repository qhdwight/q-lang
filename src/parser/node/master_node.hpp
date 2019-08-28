#pragma once

#include "abstract_node.hpp"

namespace ql::parser {
    class MasterNode : public AbstractNode {
    public:
        MasterNode() : AbstractNode(ParentRef()) {}
    };
}