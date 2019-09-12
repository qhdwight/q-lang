#pragma once

#include <parser/node/structure/parse_with_descriptor_node.hpp>


namespace ql::parser {
    class PackageNode : public ParseWithDescriptorNode {
    private:
        std::string m_Name;
    public:
        using ParseWithDescriptorNode::ParseWithDescriptorNode;

        void parse() override;
    };
}